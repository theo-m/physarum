package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
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
// (cd js && ./node_modules/webpack-cli/bin/cli.js && mv bundle.js ../public)

func main() {
	tmpl := template.Must(template.ParseFiles("physarum.gohtml"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf(r.RequestURI)
		cfg := physarum.RandomConfig()
		_ = tmpl.Execute(w, cfg)
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
				go newLaunch(c, &cfg)
			}
		}(conn)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panicf("couldnot start server: %v", err)
	}
}

func newLaunch(ws *websocket.Conn, cfg *pb.Config) {
	// TODO: while we have no input for interaction matrix
	itcMtx := physarum.RandomAttractionTable(len(cfg.Agents))
	//itcMtx := make([][]float32, len(cfg.Agents))
	//for i := 0; i < len(cfg.Agents); i++ {
	//	itcMtx[i] = make([]float32, len(cfg.Agents))
	//for j := 0; j < len(cfg.Agents); j++ {
	//	itcMtx[i][j] = cfg.InteractionMatrix[i*len(cfg.Agents)+j]
	//}
	//}
	model := physarum.NewModel(
		int(cfg.Width),
		int(cfg.Height),
		int(cfg.Particles),
		int(cfg.BlurRadius),
		int(cfg.BlurPasses),
		cfg.ZoomFactor,
		cfg.Agents,
		itcMtx,
	)
	now := time.Now().UTC().UnixNano() / 1000
	dir := fmt.Sprintf("out/%d", now)
	if err := os.MkdirAll(dir, 0777); err != nil {
		log.Panicf("couldn't create directory [%s]: %v", dir, err)
	}
	for i := 0; i < int(cfg.Iterations); i++ {
		model.Step()
		im := physarum.Image(model.W, model.H, model.Data(), model.Palette(), 0, 0, 1/2.2)
		path := fmt.Sprintf("out/%d/%d.png", now, i)
		log.Printf("saving locally iter %d", i)
		if err := physarum.SavePNG(path, im, png.DefaultCompression); err != nil {
			log.Printf("local save failed:s %v", err)
		}

		go func(ws *websocket.Conn, im *image.Image, i int) {
			var encoder png.Encoder
			encoder.CompressionLevel = png.BestCompression
			var bb bytes.Buffer
			_ = encoder.Encode(&bb, *im)
			log.Printf("sending event %d of %d", i, cfg.Iterations)

			sb, _ := proto.Marshal(&pb.Event{Content: &pb.Event_Step{Step: fmt.Sprintf("%d / %d", i, cfg.Iterations)}})
			if err := ws.WriteMessage(websocket.BinaryMessage, sb); err != nil {
				log.Printf("failed to send message: %v", err)
			}

			b, _ := proto.Marshal(&pb.Event{Content: &pb.Event_Picture{Picture: bb.Bytes()}})
			if err := ws.WriteMessage(websocket.BinaryMessage, b); err != nil {
				log.Printf("failed to send message: %v", err)
			}
		}(ws, &im, i)
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

	log.Println("sending video")
	b, err := ioutil.ReadFile(vid)
	bb, _ := proto.Marshal(&pb.Event{Content: &pb.Event_Video{Video: b}})
	if err = ws.WriteMessage(websocket.BinaryMessage, bb); err != nil {
		log.Printf("error sending video: %v", err)
	}
}
