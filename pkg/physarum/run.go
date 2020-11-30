package physarum

import (
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	sz         = 512
	width      = sz
	height     = sz
	particles  = 1 << 18
	iterations = 500
	blurRadius = 1
	blurPasses = 1
	zoomFactor = 1
)

func one(model *Model, iterations int) {
	now := time.Now().UTC().UnixNano() / 1000
	dir := fmt.Sprintf("out/%d", now)
	os.MkdirAll(dir, 0777)
	fmt.Println()
	fmt.Println(dir)
	fmt.Println(len(model.Particles), "particles")
	PrintConfigs(model.AgentConfigs, model.AttractionTable)
	SummarizeConfigs(model.AgentConfigs)
	palette := RandomPalette()
	for i := 0; i < iterations; i++ {
		model.Step()
		im := Image(model.W, model.H, model.Data(), palette, 0, 0, 1/2.2)
		path := fmt.Sprintf("out/%d/%d.png", now, i)
		if err := SavePNG(path, im, png.DefaultCompression); err != nil {
			log.Panic(err)
		}
		fmt.Printf("%d / %d\r", i, iterations)
	}
	cmd := exec.Command(
		"ffmpeg",
		"-start_number", "1",
		"-i", fmt.Sprintf("%s/%s.png", dir, "%d"),
		"-c:v", "libx264",
		"-r", "30",
		"-pix_fmt", "yuv420p",
		fmt.Sprintf("out/%d.mp4", now),
	)
	println(cmd.String())
	out, err := cmd.CombinedOutput()
	println(string(out))
	if err != nil {
		log.Panicf("failed to mp4: %v", err)
	}
}

func frames(model *Model, rate int) {
	palette := RandomPalette()

	saveImage := func(path string, w, h int, grids [][]float32, ch chan bool) {
		max := particles / float32(width*height) * 20
		im := Image(w, h, grids, palette, 0, max, 1/2.2)
		SavePNG(path, im, png.BestSpeed)
		if ch != nil {
			ch <- true
		}
	}

	ch := make(chan bool, 1)
	ch <- true
	for i := 0; ; i++ {
		fmt.Println(i)
		model.Step()
		if i%rate == 0 {
			<-ch
			path := fmt.Sprintf("frame%08d.png", i/rate)
			go saveImage(path, model.W, model.H, model.Data(), ch)
		}
	}
}

func Run() {
	if false {
		n := 2 + rand.Intn(4)
		configs := RandomAgentConfigs(n)
		table := RandomAttractionTable(n)
		model := NewModel(
			width, height, particles, blurRadius, blurPasses, zoomFactor,
			configs, table)
		frames(model, 3)
	}

	n := 2 + rand.Intn(4)
	configs := RandomAgentConfigs(n)
	table := RandomAttractionTable(n)
	model := NewModel(
		width, height, particles, blurRadius, blurPasses, zoomFactor,
		configs, table)
	start := time.Now()
	one(model, iterations)
	fmt.Println(time.Since(start))
}
