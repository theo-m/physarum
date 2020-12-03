import random
import math

import torch
import torch.nn.functional as F
from torchvision.transforms.functional import to_pil_image
from torchvision.transforms import GaussianBlur
import colour

# pip install torch torchvision tornado protobuf mypy-protobuf colour
from .physarum_pb2 import AgentConfig, Config


# seed = random.randint(0, 10_000)
seed = 1009
print(f"seed: {seed}")
torch.manual_seed(seed)


def randomCfg() -> Config:
    colors = [
        "#68bee5",
        "#ac96ea",
        "#ef6def",
        "#c65f7d",
        "#9c510b",
        "#567084",
        "#0f8efd",
        "#8a53e8",
    ]

    cfg = Config()
    cfg.width = 512
    cfg.height = 512
    n_species = random.randint(2, 8)
    cfg.particles = n_species * 100_000
    cfg.iterations = 100
    cfg.blur_radius = 3
    cfg.blur_passes = 1
    cfg.zoom_factor = 1
    cfg.gamma = 0.454
    cfg.gaussian_sigma = 2

    def unif(a, b):
        return a + (b - a) * torch.rand((n_species,))

    def norm(m, s):
        return m + s * torch.randn((n_species,))

    for sa, sd, ra, std, da, df, col in zip(
        unif(0, 2 / 3 * math.pi),
        unif(0, 64),
        unif(0, 2 / 3 * math.pi),
        # norm(math.pi / 3, math.pi / 6),
        unif(0.2, 2),
        unif(4, 6),
        unif(0.1, 0.5),
        [colour.Color(hex=h) for h in random.sample(colors, n_species)],
    ):
        sp = AgentConfig()
        sp.sensor_angle = sa
        sp.sensor_distance = sd
        sp.rotation_angle = ra
        sp.step_distance = std
        sp.deposition_amount = da
        sp.decay_factor = df
        sp.color = col.get_hex()
        cfg.agents.append(sp)

    m = -torch.ones(n_species, n_species) + 2 * torch.eye(n_species) + torch.randn((n_species, n_species))
    cfg.interaction_matrix.extend(m.view(-1).tolist())
    cfg.idist = Config.InitDistribution.CENTROIDS

    return cfg


class Physarum:
    cfg: Config
    device: str

    s: int  # nb of species
    n: int  # nb of particles
    c: int  # n // s
    particles: torch.Tensor  # [n, 3]
    grids: torch.Tensor  # [s, w, h]
    interact: torch.Tensor  # [s, s]

    # species helper (particles is chunked in contiguous bits, i.e. [0:100] is species 1, etc.)
    sensor_angles: torch.Tensor  # [n]
    sensor_distances: torch.Tensor  # [n]
    rotation_angles: torch.Tensor  # [n]
    step_distances: torch.Tensor  # [n]
    deposits: torch.Tensor  # [s]
    decay_factors: torch.Tensor  # [s]
    idxs: torch.Tensor  # [n]

    blur_kernel: torch.Tensor

    @classmethod
    def from_config(cls, cfg: Config, device: str):
        inst = cls()
        inst.cfg = cfg
        inst.device = device
        # todo: change config to avoid assert, ask directly pop size
        n, s = cfg.particles, len(cfg.agents)
        w, h = cfg.width, cfg.height
        c = n // s

        assert len(cfg.interaction_matrix) == s ** 2
        assert n % s == 0

        inst.interact = torch.tensor(cfg.interaction_matrix, device=device).view(s, s)

        inst.grids = torch.zeros(s, h, w, device=device)
        if cfg.idist == Config.InitDistribution.UNIFORM:
            inst.particles = torch.rand((n, 3), device=device) * torch.tensor((h, w, 2 * math.pi), device=device)
        elif cfg.idist == Config.InitDistribution.CENTRE:
            xys = torch.tensor([[h / 2, w / 2]]) + torch.randn((n, 2)) * min(h, w) // 8
            inst.particles = torch.cat((xys, 2 * math.pi * torch.rand((n, 1))), 1).to(device)
        elif cfg.idist == Config.InitDistribution.CENTROIDS:
            cc = []
            for _ in range(s):
                ctr = torch.tensor([[h, w]]) * torch.rand((1, 2))
                xys = ctr + torch.randn((c, 2)) * min(h, w) // random.randint(8, 16)
                ps = torch.cat((xys, 2 * math.pi * torch.rand((c, 1))), 1)
                cc.append(ps)
            inst.particles = torch.cat(cc, 0).to(device)

        def build_tensor(attr):
            ts = [torch.tensor([getattr(a, attr)] * c, device=device) for a in cfg.agents]
            return torch.cat(ts)

        inst.sensor_angles = build_tensor("sensor_angle")
        inst.sensor_distances = build_tensor("sensor_distance") * cfg.zoom_factor
        inst.rotation_angles = build_tensor("rotation_angle")
        inst.step_distances = build_tensor("step_distance") * cfg.zoom_factor
        inst.idxs = torch.arange(n, device=device) // c

        inst.deposits = torch.tensor([a.deposition_amount for a in cfg.agents], device=device)
        inst.decay_factors = torch.tensor([a.decay_factor for a in cfg.agents], device=device)

        r = 2 * cfg.blur_radius + 1
        if cfg.WhichOneof("decay_type") == "gaussian_sigma":
            gaussian1d = (
                (-((torch.arange(r, device=device).float() - r // 2) ** 2) / cfg.gaussian_sigma).exp().view(-1, 1)
            )
            gaussian2d = gaussian1d.mm(gaussian1d.t())
            inst.blur_kernel = gaussian2d.view(1, 1, r, r).repeat(s, 1, 1, 1)
        elif cfg.WhichOneof("decay_type") == "box":
            inst.blur_kernel = torch.ones(s, 1, r, r, device=device) / (r ** 2)

        inst.c = c
        inst.n = n
        inst.s = s

        return inst

    def step(self):
        # 1. combine grids: sum grids together weighted by the interaction matrix
        comb_grids = torch.einsum("ijk, im -> mjk", self.grids, self.interact)

        # 2. move particles
        x, y, a = self.particles.t()
        cx = (x + a.cos() * self.sensor_distances).to(int).clamp(0, self.cfg.height - 1)
        cy = (y + a.sin() * self.sensor_distances).to(int).clamp(0, self.cfg.width - 1)
        center_move = torch.stack((cx, cy))  # [n, 2]
        lx = (x + (a - self.sensor_angles).cos() * self.sensor_distances).to(int).clamp(0, self.cfg.height - 1)
        ly = (y + (a - self.sensor_angles).sin() * self.sensor_distances).to(int).clamp(0, self.cfg.width - 1)
        left_move = torch.stack((lx, ly))  # [n, 2]
        rx = (x + (a + self.sensor_angles).cos() * self.sensor_distances).to(int).clamp(0, self.cfg.height - 1)
        ry = (y + (a + self.sensor_angles).sin() * self.sensor_distances).to(int).clamp(0, self.cfg.width - 1)
        right_move = torch.stack((rx, ry))  # [n, 2]

        scores = []
        for moveset in [center_move, left_move, right_move]:
            i, j, k = torch.cat((self.idxs.view(1, -1), moveset), 0)
            s = comb_grids[i, j, k]
            scores.append(s)
        c, l, r = scores
        directions = self.rotation_angles * torch.randint(0, 2, (self.n,), device=self.device)
        directions[(c > l).bitwise_and(c > r)] = 0
        directions[l < r] = 1
        directions[r < l] = -1

        new_a = a + directions
        self.particles = torch.stack(
            (
                (x + new_a.cos() * self.step_distances).clamp(0, self.cfg.height - 1),
                (y + new_a.sin() * self.step_distances).clamp(0, self.cfg.width - 1),
                new_a,
            )
        ).t()

        # 3. deposit, decay + blur
        g, x, y, _ = torch.cat((self.idxs.view(-1, 1), self.particles.to(int)), 1).t()
        self.grids[g, x, y] += self.deposits[g]
        self.grids = F.conv2d(
            self.grids.unsqueeze(0),
            self.blur_kernel,
            stride=1,
            groups=self.s,
            padding=self.cfg.blur_radius,
        ).squeeze(0)
        self.grids *= self.decay_factors.view(-1, 1, 1)

    def img(self) -> torch.Tensor:
        maxs = torch.tensor([m.quantile(0.99) for m in self.grids], device=self.device).view(-1, 1, 1)
        r = torch.zeros((self.cfg.height, self.cfg.width), device=self.device)
        g, b = r.clone(), r.clone()
        normGrids = (self.grids / maxs).clamp(0, 1) ** self.cfg.gamma
        for i, a in enumerate(self.cfg.agents):
            c = colour.Color()
            c.set_hex(a.color)
            r += int(255 * c.get_red()) * normGrids[i]
            g += int(255 * c.get_green()) * normGrids[i]
            b += int(255 * c.get_blue()) * normGrids[i]
        img = torch.stack(
            (
                r.clamp(0, 255).to(torch.uint8),
                g.clamp(0, 255).to(torch.uint8),
                b.clamp(0, 255).to(torch.uint8),
            )
        )
        return img


if __name__ == "__main__":
    import time
    import subprocess
    from pathlib import Path

    t = int(time.time())
    dirr = Path(f"out/{t}")
    dirr.mkdir(parents=True, exist_ok=False)

    cfg = randomCfg()
    p = Physarum.from_config(cfg, "cpu")
    print(p.cfg)
    for i in range(cfg.iterations):
        print(f"\r{i:5d}/{cfg.iterations}", end="")
        p.step()
        imt = p.img()
        img = to_pil_image(imt)
        img.save(f"{dirr}/{i}.png")

    cmd = f"ffmpeg -start_number 1 -i {dirr}/%d.png -c:v libx264 -r 30 -pix_fmt yuv420p out/{t}.mp4"
    subprocess.run(cmd.split(" "))
