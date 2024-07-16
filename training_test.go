package som

import (
	"testing"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/neighborhood"
	"github.com/stretchr/testify/assert"
)

func TestNewTrainer(t *testing.T) {
	params := TrainingParams{}
	somParams := SomParams{
		Size: Size{2, 3},
		Layers: []LayerDef{
			{
				Columns: []string{"x", "y"},
				Weight:  0.5,
			},
			{
				Columns: []string{"a", "b", "c"},
				Weight:  1.0,
			},
		},
	}
	som := New(&somParams)

	t1 := NewTable([]string{"x", "y", "a", "b", "c"}, 5)
	_, err := NewTrainer(&som, &t1, &params)
	assert.Nil(t, err)

	t1 = NewTable([]string{"x", "y", "a", "b"}, 5)
	_, err = NewTrainer(&som, &t1, &params)
	assert.NotNil(t, err)

	t1 = NewTable([]string{"x", "y", "a", "b", "c", "d"}, 5)
	_, err = NewTrainer(&som, &t1, &params)
	assert.NotNil(t, err)

	t1 = NewTable([]string{"x", "a", "y", "b", "c"}, 5)
	_, err = NewTrainer(&som, &t1, &params)
	assert.NotNil(t, err)
}

func TestTrainerTrain(t *testing.T) {
	params := TrainingParams{
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Power{Start: 5, End: 0.5},
	}
	somParams := SomParams{
		Size: Size{2, 3},
		Layers: []LayerDef{
			{
				Columns: []string{"x", "y"},
				Weight:  0.5,
			},
			{
				Columns: []string{"a", "b", "c"},
				Weight:  1.0,
			},
		},
		Neighborhood: &neighborhood.Gaussian{},
	}
	som := New(&somParams)

	t1 := NewTable([]string{"x", "y", "a", "b", "c"}, 5)
	trainer, err := NewTrainer(&som, &t1, &params)
	assert.Nil(t, err)

	trainer.Train(10)
}
