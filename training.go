package som

import (
	"math/rand"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/table"
)

type TrainingConfig struct {
	Epochs             int
	LearningRate       decay.Decay
	NeighborhoodRadius decay.Decay
	ViSomLambda        float64
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

func (t *Trainer) Train(progress chan float64) {
	t.som.randomize(t.rng)

	var meanDist float64
	for epoch := 0; epoch < t.params.Epochs; epoch++ {
		meanDist = t.epoch(epoch)
		progress <- meanDist
	}
	progress <- meanDist

	close(progress)
}

func (t *Trainer) epoch(epoch int) float64 {
	alpha := t.params.LearningRate.Decay(epoch, t.params.Epochs)
	radius := t.params.NeighborhoodRadius.Decay(epoch, t.params.Epochs)

	data := make([][]float64, len(t.tables))
	rows := t.tables[0].Rows()

	dist := 0.0
	for i := 0; i < rows; i++ {
		for j := 0; j < len(t.tables); j++ {
			data[j] = t.tables[j].GetRow(i)
		}
		dist += t.som.learn(data, alpha, radius, t.params.ViSomLambda)

		if t.params.ViSomLambda == 0 || i%10 != 0 { // SOM
			continue
		}
		// ViSOM refresh: present random node as data
		node := t.rng.Intn(t.som.size.Nodes())
		for j := 0; j < len(t.tables); j++ {
			data[j] = t.som.layers[j].GetNodeAt(node)
		}
		t.som.learn(data, alpha, radius, t.params.ViSomLambda)
	}

	return dist / float64(rows)
}
