package som

import (
	"fmt"
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
	t := &Trainer{
		som:    som,
		tables: tables,
		params: params,
		rng:    rng,
	}
	if err := t.checkTable(); err != nil {
		return nil, err
	}
	return t, nil
}

// checkTable checks that the table columns match the SOM layer columns.
// It returns true if the table columns match the SOM layer columns, false otherwise.
func (t *Trainer) checkTable() error {
	if len(t.tables) == 0 {
		return fmt.Errorf("no tables provided")
	}

	if len(t.som.layers) != len(t.tables) {
		return fmt.Errorf("number of tables (%d) does not match number of layers (%d)", len(t.tables), len(t.som.layers))
	}

	rows := -1
	for _, table := range t.tables {
		if rows == -1 {
			rows = table.Rows()
		} else if rows != table.Rows() {
			return fmt.Errorf("number of rows in table (%d) does not match number of rows in table (%d)", rows, table.Rows())
		}
	}

	for i := range t.som.layers {
		table := t.tables[i]
		cols := t.som.layers[i].columns
		if table.Columns() != len(cols) {
			return fmt.Errorf("number of columns in table (%d) does not match number of columns in layer (%d)", table.Columns(), len(cols))
		}
		for j, col := range cols {
			if table.ColumnNames()[j] != col {
				return fmt.Errorf("column %d in table (%s) does not match column %d in layer (%s)", j, table.ColumnNames()[j], j, col)
			}
		}
	}
	return nil
}

func (t *Trainer) Train(maxEpoch int) {
	t.randomize()
	for epoch := 0; epoch < maxEpoch; epoch++ {
		t.epoch(epoch, maxEpoch)
	}
}

func (t *Trainer) randomize() {
	for _, layer := range t.som.layers {
		for i := range layer.data {
			layer.data[i] = t.rng.Float64()
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
