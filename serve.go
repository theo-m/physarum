package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/theo-m/physarum/pkg/pb"
	"github.com/theo-m/physarum/pkg/physarum"
)

//go:generate protoc --go_out=pkg/pb --go_opt=paths=source_relative physarum.proto
//go:generate protoc --js_out=import_style=commonjs,binary:js physarum.proto
// (cd js && ./node_modules/webpack-cli/bin/cli.js && cp bundle.js ../public)

func main() {
	tmpl := template.Must(template.ParseFiles("physarum.gohtml"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf(r.RequestURI)
		tmpl.Execute(w, nil)
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upd := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
		conn, err := upd.Upgrade(w, r, nil)
		if err != nil {
			log.Panicf("failed to build ws connection: %v", err)
		}
		log.Printf("websocket conn")
		go func(c *websocket.Conn) {
			defer c.Close()
			for {
				_, content, err := c.ReadMessage()
				if err != nil {
					log.Printf("could not read message: %v", err)
					return
				}
				var cfg pb.Config
				if err = proto.Unmarshal(content, &cfg); err != nil {
					log.Println(string(content))
					log.Printf("failed to parse form: %v", err)
					return
				}
				go newLaunch(c, cfg)
			}
		}(conn)
	})
	http.ListenAndServe(":8080", nil)
}

func newLaunch(ws *websocket.Conn, cfg pb.Config) {
	model := physarum.NewModel(
		int(cfg.Width),
		int(cfg.Height),
		int(cfg.Particles),
		int(cfg.BlurRadius),
		int(cfg.BlurPasses),
		cfg.ZoomFactor,
		physarum.RandomConfigs(int(cfg.NumGrids)),
		physarum.RandomAttractionTable(int(cfg.NumGrids)),
	)
	palette := physarum.RandomPalette()
	now := time.Now().UTC().UnixNano() / 1000
	dir := fmt.Sprintf("out/%d", now)
	os.MkdirAll(dir, 0777)
	for i := 0; i < int(cfg.Iterations); i++ {
		model.Step()
		im := physarum.Image(model.W, model.H, model.Data(), palette, 0, 0, 1/2.2)
		path := fmt.Sprintf("out/%d/%d.png", now, i)
		if err := physarum.SavePNG(path, im, png.DefaultCompression); err != nil {
			log.Printf("local save failed:s %v", err)
		}

		var encoder png.Encoder
		encoder.CompressionLevel = png.BestCompression
		var bb bytes.Buffer
		encoder.Encode(&bb, im)
		log.Printf("iter %d of %d", i, cfg.Iterations)

		sb, _ := proto.Marshal(&pb.Event{Content: &pb.Event_Step{Step: fmt.Sprintf("%d / %d", i, cfg.Iterations)}})
		if err := ws.WriteMessage(websocket.BinaryMessage, sb); err != nil {
			log.Printf("failed to send message: %v", err)
		}

		b, _ := proto.Marshal(&pb.Event{Content: &pb.Event_Picture{Picture: bb.Bytes()}})
		if err := ws.WriteMessage(websocket.BinaryMessage, b); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
	vid := fmt.Sprintf("out/%d.mp4", now)
	cmd := exec.Command(
		"ffmpeg",
		"-start_number", "1",
		"-i", fmt.Sprintf("%s/%s.png", dir, "%d"),
		"-c:v", "libx264",
		"-r", "30",
		"-pix_fmt", "yuv420p",
		vid,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Panicf("failed to mp4: %v", err)
	}
	b, err := ioutil.ReadFile(vid)
	log.Println("sending video")
	ev := &pb.Event{Content: &pb.Event_Video{Video: b}}
	bb, _ := proto.Marshal(ev)
	ws.WriteMessage(websocket.BinaryMessage, bb)
}
