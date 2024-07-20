package som

import (
	"math/rand"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/table"
)

type TrainingConfig struct {
	LearningRate       decay.Decay
	NeighborhoodRadius decay.Decay
	VisualLambda       float64
}

type Trainer struct {
	som    *Som
	tables []*table.Table
	params *TrainingConfig
	rng    *rand.Rand
}

func NewTrainer(som *Som, tables []*table.Table, params *TrainingConfig, rng *rand.Rand) (*Trainer, error) {
	if err := checkTables(som, tables); err != nil {
		return nil, err
	}
	return &Trainer{
		som:    som,
		tables: tables,
		params: params,
		rng:    rng,
	}, nil
}

func (t *Trainer) Train(maxEpoch int, progress chan float64) {
	t.som.randomize(t.rng)

	var meanDist float64
	for epoch := 0; epoch < maxEpoch; epoch++ {
		meanDist = t.epoch(epoch, maxEpoch)
		progress <- meanDist
	}
	progress <- meanDist

	close(progress)
}

func (t *Trainer) epoch(epoch, maxEpoch int) float64 {
	alpha := t.params.LearningRate.Decay(epoch, maxEpoch)
	radius := t.params.NeighborhoodRadius.Decay(epoch, maxEpoch)

	data := make([][]float64, len(t.tables))
	rows := t.tables[0].Rows()

	dist := 0.0
	for i := 0; i < rows; i++ {
		for j := 0; j < len(t.tables); j++ {
			data[j] = t.tables[j].GetRow(i)
		}
		dist += t.som.learn(data, alpha, radius, t.params.VisualLambda)

		if t.params.VisualLambda == 0 || i%10 != 0 { // SOM
			continue
		}
		// ViSOM refresh: present random node as data
		node := t.rng.Intn(t.som.size.Width * t.som.size.Height)
		for j := 0; j < len(t.tables); j++ {
			data[j] = t.som.layers[j].GetNodeAt(node)
		}
		t.som.learn(data, alpha, radius, t.params.VisualLambda)
	}

	return dist / float64(rows)
}
