import logging
import time
import subprocess
from io import BytesIO
from typing import Union, Optional, Awaitable
from pathlib import Path

import torch
import tornado.ioloop
import tornado.web
import tornado.websocket
import tornado.log
import tornado.template
from torchvision.transforms.functional import to_pil_image
from PIL import Image

from .physarum_pb2 import Event, Config
from .model import randomCfg, Physarum

logging.basicConfig(level=logging.DEBUG)
log = logging.getLogger()
[h.setFormatter(tornado.log.LogFormatter()) for h in log.handlers]


def configpath(name: str) -> Path:
    return f"{name}-cfg.pb"


def videopath(name: str) -> Path:
    return f"{name}.mp4"


class MainHandler(tornado.web.RequestHandler):
    def get(self):
        from_config = self.get_query_argument("from-config", None)
        cfgp = Path(f"out/{configpath(from_config)}")
        if cfgp.exists() and cfgp.is_file():
            cfg = Config()
            with cfgp.open("rb") as fi:
                cfg.ParseFromString(fi.read())
            video = f"/sims/{videopath(from_config)}"
        else:
            cfg = randomCfg()
            video = ""
        self.render("physarum.html", cfg=cfg, video=video)


class ListHandler(tornado.web.RequestHandler):
    def get(self):
        files = Path("out").iterdir()
        sims_dir = [f for f in files if f.is_dir()]
        sims_names = [f.name for f in sims_dir if Path(f"out/{configpath(f.name)}").exists()]
        last_imgs = [sorted(Path(f"out/{n}").iterdir(), key=lambda p: p.lstat().st_atime)[-1].name for n in sims_names]
        self.render("lists.html", sims=[(n, f"/sims/{n}/{im}") for n, im in zip(sims_names, last_imgs)])


class WSHandler(tornado.websocket.WebSocketHandler):
    def initialize(self, device: str):
        self.device = device

    def open(self):
        pass

    def on_message(self, message: Union[str, bytes]) -> Optional[Awaitable[None]]:
        logging.info("received new config")
        cfg = Config()
        cfg.ParseFromString(message)

        n = len(cfg.agents)
        m = -torch.ones(n, n) + 2 * torch.eye(n) + torch.randn((n, n))
        cfg.interaction_matrix.extend(m.view(-1).tolist())

        model = Physarum.from_config(cfg, self.device)
        mname = int(time.time())

        logging.info(f"sending model id '{mname}'")
        ev = Event()
        ev.model_id = str(mname)
        self.write_message(ev.SerializeToString(), binary=True)

        mdir = Path(f"out/{mname}")
        mdir.mkdir(parents=True, exist_ok=False)

        with open(f"out/{mname}-cfg.pb", "wb") as fo:
            fo.write(model.cfg.SerializeToString())

        for it in range(cfg.iterations):
            model.step()
            imt = model.img()
            img: Image = to_pil_image(imt)
            img.save(f"{mdir}/{it}.png")

            ev = Event()
            ev.step = f"{it} / {cfg.iterations}"
            self.write_message(ev.SerializeToString(), binary=True)

            if it % 5 == 0:
                ev = Event()
                bytes = BytesIO()
                img.save(bytes, format="PNG", optimize=True)
                ev.picture = bytes.getvalue()
                self.write_message(ev.SerializeToString(), binary=True)

        vidp = Path(f"out/{mname}.mp4")
        cmd = f"ffmpeg -start_number 1 -i {mdir}/%d.png -c:v libx264 -r 30 -pix_fmt yuv420p {vidp}"
        subprocess.run(cmd.split(" "), capture_output=True)

        ev = Event()
        with vidp.open("rb") as fi:
            vid_bytes = fi.read()
        ev.video = vid_bytes
        self.write_message(ev.SerializeToString(), binary=True)
        return


def make_app(device: str):
    return tornado.web.Application(
        [
            (r"/", MainHandler),
            (r"/list", ListHandler),
            (r"/ws", WSHandler, {"device": device}),
            (r"/public/(.*)", tornado.web.StaticFileHandler, {"path": "./public"}),
            (r"/sims/(.*)", tornado.web.StaticFileHandler, {"path": "./out"}),
        ]
    )


if __name__ == "__main__":
    import argparse

    ap = argparse.ArgumentParser()
    ap.add_argument("-p", type=int, default=8080)
    ap.add_argument("-d", choices=["cpu", "cuda"], default="cpu")
    args = ap.parse_args()

    app = make_app(args.d)
    app.listen(args.p)
    tornado.ioloop.IOLoop.current().start()
