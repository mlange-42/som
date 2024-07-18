package som

import (
	"math/rand"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/table"
)

type TrainingConfig struct {
	LearningRate       decay.Decay
	NeighborhoodRadius decay.Decay
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

func (t *Trainer) Train(maxEpoch int) {
	t.randomize()
	for epoch := 0; epoch < maxEpoch; epoch++ {
		t.epoch(epoch, maxEpoch)
	}
}

func (t *Trainer) randomize() {
	for _, lay := range t.som.layers {
		data := lay.Data()
		for i := range data {
			data[i] = t.rng.Float64()
		}
	}
}

func (t *Trainer) epoch(epoch, maxEpoch int) {
	alpha := t.params.LearningRate.Decay(epoch, maxEpoch)
	radius := t.params.NeighborhoodRadius.Decay(epoch, maxEpoch)

	//fmt.Println("Epoch", epoch, "of", maxEpoch, "alpha", alpha, "radius", radius)

	data := make([][]float64, len(t.tables))
	rows := t.tables[0].Rows()

	for i := 0; i < rows; i++ {
		for j := 0; j < len(t.tables); j++ {
			data[j] = t.tables[j].GetRow(i)
		}
		t.som.learn(data, alpha, radius)
	}
}
