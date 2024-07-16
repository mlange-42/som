package som

import (
	"fmt"

	"github.com/mlange-42/som/decay"
)

type TrainingParams struct {
	LearningRate       decay.Decay
	NeighborhoodRadius decay.Decay
}

type Trainer struct {
	som    *Som
	table  *Table
	params *TrainingParams
}

func NewTrainer(som *Som, table *Table, params *TrainingParams) (*Trainer, error) {
	t := &Trainer{
		som:    som,
		table:  table,
		params: params,
	}
	if !t.checkTable() {
		return nil, fmt.Errorf("table columns do not match SOM layer columns")
	}
	return t, nil
}

// checkTable checks that the table columns match the SOM layer columns.
// It returns true if the table columns match the SOM layer columns, false otherwise.
func (t *Trainer) checkTable() bool {
	if len(t.table.columns) != t.som.offset[len(t.som.offset)-1]+len(t.som.layers[len(t.som.layers)-1].columns) {
		return false
	}
	for i := range t.som.layers {
		off := t.som.offset[i]
		for j, col := range t.som.layers[i].columns {
			if t.table.columns[j+off] != col {
				return false
			}
		}
	}
	return true
}

func (t *Trainer) Train(maxEpoch int) {
	for epoch := 0; epoch < maxEpoch; epoch++ {
		t.epoch(epoch, maxEpoch)
	}
}

func (t *Trainer) epoch(epoch, maxEpoch int) {
	alpha := t.params.LearningRate.Decay(epoch, maxEpoch-1)
	radius := t.params.NeighborhoodRadius.Decay(epoch, maxEpoch-1)

	fmt.Println("Epoch", epoch, "of", maxEpoch, "alpha", alpha, "radius", radius)

	for i := range t.table.rows {
		t.som.learn(t.table.GetRow(i), alpha, radius)
	}
}
