package physarum

import (
	"math/rand"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/theo-m/physarum/pkg/pb"
)

const (
	sensorAngleMin      = 0
	sensorAngleMax      = 120
	sensorDistanceMin   = 0
	sensorDistanceMax   = 64
	rotationAngleMin    = 0
	rotationAngleMax    = 120
	stepDistanceMin     = 0.2
	stepDistanceMax     = 2
	depositionAmountMin = 4
	depositionAmountMax = 6
	decayFactorMin      = 0.1
	decayFactorMax      = 0.5

	attractionFactorMean = 1
	attractionFactorStd  = 0.5
	repulsionFactorMean  = -1
	repulsionFactorStd   = 0.5
)

func RandomAgentConfig() *pb.AgentConfig {
	uniform := func(min, max float32) float32 {
		return min + rand.Float32()*(max-min)
	}

	sensorAngle := Radians(uniform(sensorAngleMin, sensorAngleMax))
	sensorDistance := uniform(sensorDistanceMin, sensorDistanceMax)
	rotationAngle := Radians(uniform(rotationAngleMin, rotationAngleMax))
	stepDistance := uniform(stepDistanceMin, stepDistanceMax)
	depositionAmount := uniform(depositionAmountMin, depositionAmountMax)
	decayFactor := uniform(decayFactorMin, decayFactorMax)

	return &pb.AgentConfig{
		SensorAngle:      sensorAngle,
		SensorDistance:   sensorDistance,
		RotationAngle:    rotationAngle,
		StepDistance:     stepDistance,
		DepositionAmount: depositionAmount,
		DecayFactor:      decayFactor,
	}
}

func RandomAgentConfigs(n int) []*pb.AgentConfig {
	configs := make([]*pb.AgentConfig, n)
	palette, _ := colorful.HappyPalette(n)
	for i := range configs {
		configs[i] = RandomAgentConfig()
		configs[i].Color = palette[i].Hex()
	}
	return configs
}

func RandomAttractionTable(n int) [][]float32 {
	normal := func(mean, std float32) float32 {
		return mean + float32(rand.NormFloat64())*std
	}

	result := make([][]float32, n)
	for i := range result {
		result[i] = make([]float32, n)
		for j := range result[i] {
			if i == j {
				result[i][j] = normal(attractionFactorMean, attractionFactorStd)
			} else {
				result[i][j] = normal(repulsionFactorMean, repulsionFactorStd)
			}
		}
	}
	return result
}

func RandomConfig() *pb.Config {
	n := 4 + rand.Intn(4)
	itcMtx := make([]float32, n*n)
	rndMtx := RandomAttractionTable(n)
	for i := 0; i < n*n; i++ {
		itcMtx[i] = rndMtx[i/n][i%n]
	}
	return &pb.Config{
		Width:             512,
		Height:            512,
		Particles:         1 << (11 + rand.Intn(10)),
		Iterations:        (1 + rand.Int31n(20)) * 100,
		BlurRadius:        2,
		BlurPasses:        1,
		ZoomFactor:        1,
		Agents:            RandomAgentConfigs(n),
		InteractionMatrix: itcMtx,
	}
}
