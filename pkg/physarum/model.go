package physarum

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/theo-m/physarum/pkg/pb"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"sync"
)

type Model struct {
	W int
	H int

	BlurRadius int
	BlurPasses int

	ZoomFactor float32

	AgentConfigs    []*pb.AgentConfig
	AttractionTable [][]float32

	Grids     []*Grid
	Particles []Particle

	Iteration int
}

func NewModel(
	w, h, numParticles, blurRadius, blurPasses int,
	zoomFactor float32,
	configs []*pb.AgentConfig,
	attractionTable [][]float32,
	distrib pb.Config_InitDistribution,
) *Model {

	grids := make([]*Grid, len(configs))
	numParticlesPerConfig := int(math.Ceil(
		float64(numParticles) / float64(len(configs))))
	actualNumParticles := numParticlesPerConfig * len(configs)
	particles := make([]Particle, actualNumParticles)
	m := &Model{
		w,
		h,
		blurRadius,
		blurPasses,
		zoomFactor,
		configs,
		attractionTable,
		grids,
		particles,
		0,
	}
	m.StartOver(distrib)
	return m
}

func (m *Model) StartOver(distribution pb.Config_InitDistribution) {
	numParticlesPerConfig := len(m.Particles) / len(m.AgentConfigs)
	m.Particles = m.Particles[:0]
	m.Iteration = 0
	switch distribution {
	case pb.Config_UNIFORM:
		for c := range m.AgentConfigs {
			m.Grids[c] = NewGrid(m.W, m.H)
			for i := 0; i < numParticlesPerConfig; i++ {
				x := rand.Float32() * float32(m.W)
				y := rand.Float32() * float32(m.H)
				a := rand.Float32() * 2 * math.Pi
				p := Particle{x, y, a, uint32(c)}
				m.Particles = append(m.Particles, p)
			}
		}
	case pb.Config_CENTROIDS:
		ws, hs := float32(m.W)/float32(len(m.AgentConfigs)), float32(m.H)/float32(len(m.AgentConfigs))
		for c := range m.AgentConfigs {
			m.Grids[c] = NewGrid(m.W, m.H)
			pcx, pcy := rand.Float32()*2*ws, rand.Float32()*2*hs
			for i := 0; i < numParticlesPerConfig; i++ {
				x := float32(rand.NormFloat64())*ws + pcx
				y := float32(rand.NormFloat64())*hs + pcy
				a := rand.Float32() * 2 * math.Pi
				p := Particle{x, y, a, uint32(c)}
				m.Particles = append(m.Particles, p)
			}
		}
	case pb.Config_CENTRE:
		for c := range m.AgentConfigs {
			m.Grids[c] = NewGrid(m.W, m.H)
			for i := 0; i < numParticlesPerConfig; i++ {
				x := float32(rand.NormFloat64()) * float32(m.W) / 4
				y := float32(rand.NormFloat64()) * float32(m.H) / 4
				a := rand.Float32() * 2 * math.Pi
				p := Particle{x, y, a, uint32(c)}
				m.Particles = append(m.Particles, p)
			}
		}
	case pb.Config_GRID:
		offset := float32(10)
		width := float32(5)
		nlines := 10
		for c := range m.AgentConfigs {
			m.Grids[c] = NewGrid(m.W, m.H)
			for i := 0; i < numParticlesPerConfig; i++ {
				var x, y, a float32
				if c%2 == 0 { // vertical lines
					x = offset*float32(c) + rand.Float32()*width*float32(m.W)/float32(len(m.AgentConfigs))*float32(i%nlines-nlines/2)
					y = rand.Float32() * float32(m.H)
					a = math.Pi/2*float32(math.Pow(-1, float64(rand.Intn(9)%2))) + rand.Float32()/4
				} else { // horizontal lines
					x = rand.Float32() * float32(m.W)
					y = offset*float32(c) + rand.Float32()*width*float32(m.H)/float32(len(m.AgentConfigs))*float32(i%nlines-nlines/2)
					a = math.Pi*float32(math.Pow(-1, float64(rand.Intn(9)%2))) + rand.Float32()/4
				}
				p := Particle{x, y, a, uint32(c)}
				m.Particles = append(m.Particles, p)
			}
		}
	}
	ws, hs := float32(m.W)/float32(len(m.AgentConfigs)), float32(m.H)/float32(len(m.AgentConfigs))
	for c := range m.AgentConfigs {
		m.Grids[c] = NewGrid(m.W, m.H)
		//pcx, pcy := rand.Float32() * 2 * ws, rand.Float32() * 2 * hs
		for i := 0; i < numParticlesPerConfig; i++ {
			x := float32(rand.NormFloat64()) * ws
			y := float32(rand.NormFloat64()) * hs
			a := rand.Float32() * 2 * math.Pi
			p := Particle{x, y, a, uint32(c)}
			m.Particles = append(m.Particles, p)
		}
	}
}

func (m *Model) Step() {
	updateParticle := func(rnd *rand.Rand, i int) {
		p := m.Particles[i]
		config := m.AgentConfigs[p.C]
		grid := m.Grids[p.C]

		// u := p.X / float32(m.W)
		// v := p.Y / float32(m.H)

		sensorDistance := config.SensorDistance * m.ZoomFactor
		sensorAngle := config.SensorAngle
		rotationAngle := config.RotationAngle
		stepDistance := config.StepDistance * m.ZoomFactor

		xc := p.X + cos(p.A)*sensorDistance
		yc := p.Y + sin(p.A)*sensorDistance
		xl := p.X + cos(p.A-sensorAngle)*sensorDistance
		yl := p.Y + sin(p.A-sensorAngle)*sensorDistance
		xr := p.X + cos(p.A+sensorAngle)*sensorDistance
		yr := p.Y + sin(p.A+sensorAngle)*sensorDistance
		C := grid.GetTemp(xc, yc)
		L := grid.GetTemp(xl, yl)
		R := grid.GetTemp(xr, yr)

		da := rotationAngle * direction(rnd, C, L, R)
		// da := rotationAngle * weightedDirection(rnd, C, L, R)
		p.A = Shift(p.A+da, 2*math.Pi)
		p.X = Shift(p.X+cos(p.A)*stepDistance, float32(m.W))
		p.Y = Shift(p.Y+sin(p.A)*stepDistance, float32(m.H))
		m.Particles[i] = p
	}

	updateParticles := func(wi, wn int, wg *sync.WaitGroup) {
		seed := int64(m.Iteration)<<8 | int64(wi)
		rnd := rand.New(rand.NewSource(seed))
		n := len(m.Particles)
		batch := int(math.Ceil(float64(n) / float64(wn)))
		i0 := wi * batch
		i1 := i0 + batch
		if wi == wn-1 {
			i1 = n
		}
		for i := i0; i < i1; i++ {
			updateParticle(rnd, i)
		}
		wg.Done()
	}

	updateGrids := func(c int, wg *sync.WaitGroup) {
		config := m.AgentConfigs[c]
		grid := m.Grids[c]
		for _, p := range m.Particles {
			if uint32(c) == p.C {
				grid.Add(p.X, p.Y, config.DepositionAmount)
			}
		}
		grid.BoxBlur(m.BlurRadius, m.BlurPasses, config.DecayFactor)
		wg.Done()
	}

	combineGrids := func(c int, wg *sync.WaitGroup) {
		grid := m.Grids[c]
		for i := range grid.Temp {
			grid.Temp[i] = 0
		}
		for i, other := range m.Grids {
			factor := m.AttractionTable[c][i]
			for j, value := range other.Data {
				grid.Temp[j] += value * factor
			}
		}
		wg.Done()
	}

	var wg sync.WaitGroup

	// step 1: combine grids
	for i := range m.AgentConfigs {
		wg.Add(1)
		go combineGrids(i, &wg)
	}
	wg.Wait()

	// step 2: move particles
	wn := runtime.NumCPU()
	for wi := 0; wi < wn; wi++ {
		wg.Add(1)
		go updateParticles(wi, wn, &wg)
	}
	wg.Wait()

	// step 3: deposit, blur, and decay
	for i := range m.AgentConfigs {
		wg.Add(1)
		go updateGrids(i, &wg)
	}
	wg.Wait()

	m.Iteration++
}

func (m *Model) Data() [][]float32 {
	result := make([][]float32, len(m.Grids))
	for i, grid := range m.Grids {
		result[i] = make([]float32, len(grid.Data))
		copy(result[i], grid.Data)
	}
	return result
}

func (m *Model) Palette() []color.RGBA {
	p := make([]color.RGBA, len(m.AgentConfigs))
	for i, ac := range m.AgentConfigs {
		c, _ := colorful.Hex(ac.Color)
		r, g, b, a := c.RGBA()
		p[i] = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}
	return p
}

func direction(rnd *rand.Rand, C, L, R float32) float32 {
	if C > L && C > R {
		return 0
	} else if C < L && C < R {
		return float32((rnd.Int63()&1)<<1 - 1)
	} else if L < R {
		return 1
	} else if R < L {
		return -1
	}
	return 0
}

func weightedDirection(rnd *rand.Rand, C, L, R float32) float32 {
	W := [3]float32{C, L, R}
	D := [3]float32{0, -1, 1}

	if W[0] > W[1] {
		W[0], W[1] = W[1], W[0]
		D[0], D[1] = D[1], D[0]
	}
	if W[0] > W[2] {
		W[0], W[2] = W[2], W[0]
		D[0], D[2] = D[2], D[0]
	}
	if W[1] > W[2] {
		W[1], W[2] = W[2], W[1]
		D[1], D[2] = D[2], D[1]
	}

	a := W[1] - W[0]
	b := W[2] - W[1]
	if rnd.Float32()*(a+b) < a {
		return D[1]
	}
	return D[2]
}
