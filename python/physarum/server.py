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


class MainHandler(tornado.web.RequestHandler):
    def get(self):
        cfg = randomCfg()
        self.render("physarum.html", cfg=cfg)


class WSHandler(tornado.websocket.WebSocketHandler):
    def open(self):
        pass

    def on_message(self, message: Union[str, bytes]) -> Optional[Awaitable[None]]:
        logging.info("received new config")
        cfg = Config()
        cfg.ParseFromString(message)

        n = len(cfg.agents)
        m = -torch.ones(n, n) + 2 * torch.eye(n) + torch.randn((n, n))
        cfg.interaction_matrix.extend(m.view(-1).tolist())

        model = Physarum.from_config(cfg)
        mname = int(time.time())
        mdir = Path(f"out/{mname}")
        mdir.mkdir()
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


def make_app():
    return tornado.web.Application(
        [
            (r"/", MainHandler),
            (r"/ws", WSHandler),
            (r"/public/(.*)", tornado.web.StaticFileHandler, {"path": "./public"}),
        ]
    )


if __name__ == "__main__":
    app = make_app()
    app.listen(8988)
    tornado.ioloop.IOLoop.current().start()
