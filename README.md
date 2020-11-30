# Physarum Transport Networks

## Web GUI fork

Need go installed, the websocket event are serialized with protobuf so need `protoc` installed and the codegen for javascript is bad for imports so needed to run the codegen through webpack.
Run `bash build.sh` and then `docker run -p 8080:8080 <build_hash>`.

Bad gif demo (in the process of hosting the container somewhere to make it accessible): 
![1606767333](https://user-images.githubusercontent.com/17948980/100661928-7c438000-3354-11eb-99da-d7345ce63fde.gif)

Left todo: 
- upload videos to a bucket and make them browsable with config search e.g. with routing `/:hashOfConfig`.
- allow tweaking repulsion matrix (no idea how)

---

![Example](https://www.michaelfogleman.com/static/physarum/header.png)

This is a particle-based simulation inspired by the _Physarum polycephalum_
slime mold.

Details about the algorithm can be found on Sage Jenson's excellent webpage:

https://sagejenson.com/physarum

The algorithm is surprisingly simple given how complex its outputs appear.
That's the magic of generative algorithms!

## Usage

    go run cmd/physarum/main.go

## Examples

![Montage](https://www.michaelfogleman.com/static/physarum/montage-small.jpg)
