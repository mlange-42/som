package som

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/neighborhood"
	"github.com/stretchr/testify/assert"
)

func TestNewTrainer(t *testing.T) {
	params := TrainingConfig{}
	somParams := SomConfig{
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
	som, err := New(&somParams)
	assert.NoError(t, err)

	rng := rand.New(rand.NewSource(1))

	t1 := []*Table{
		NewTable([]string{"x", "y"}, 5),
		NewTable([]string{"a", "b", "c"}, 5),
	}
	_, err = NewTrainer(&som, t1, &params, rng)
	assert.Nil(t, err)

	t1 = []*Table{
		NewTable([]string{"x", "y"}, 5),
		NewTable([]string{"a", "b"}, 5),
	}
	_, err = NewTrainer(&som, t1, &params, rng)
	assert.NotNil(t, err)

	t1 = []*Table{
		NewTable([]string{"x", "y"}, 5),
		NewTable([]string{"a", "b", "c", "d"}, 5),
	}
	_, err = NewTrainer(&som, t1, &params, rng)
	assert.NotNil(t, err)

	t1 = []*Table{
		NewTable([]string{"x", "y"}, 5),
		NewTable([]string{"a", "c", "b"}, 5),
	}
	_, err = NewTrainer(&som, t1, &params, rng)
	assert.NotNil(t, err)
}

func TestTrainerDecay(t *testing.T) {
	params := TrainingConfig{
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Power{Start: 5, End: 0.5},
	}
	somParams := SomConfig{
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
	som, err := New(&somParams)
	assert.NoError(t, err)

	rng := rand.New(rand.NewSource(1))

	t1 := []*Table{
		NewTable([]string{"x", "y"}, 5),
		NewTable([]string{"a", "b", "c"}, 5),
	}
	trainer, err := NewTrainer(&som, t1, &params, rng)
	assert.Nil(t, err)

	trainer.Train(10)
}

func TestTrainerTrain(t *testing.T) {
	params := TrainingConfig{
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Power{Start: 3, End: 0.5},
	}
	somParams := SomConfig{
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
	som, err := New(&somParams)
	assert.NoError(t, err)

	rng := rand.New(rand.NewSource(1))

	t.Run("Train with zero epochs", func(t *testing.T) {
		tables := []*Table{
			NewTable([]string{"x", "y"}, 5),
			NewTable([]string{"a", "b", "c"}, 5),
		}
		trainer, err := NewTrainer(&som, tables, &params, rng)
		assert.Nil(t, err)

		trainer.Train(0)

		for _, v := range som.layers[0].data {
			assert.NotEqual(t, 0, v)
		}
	})

	t.Run("Train with one epoch", func(t *testing.T) {
		tables := []*Table{
			NewTable([]string{"x", "y"}, 5),
			NewTable([]string{"a", "b", "c"}, 5),
		}
		trainer, err := NewTrainer(&som, tables, &params, rng)
		assert.Nil(t, err)

		trainer.Train(1)
	})

	t.Run("Train with multiple epochs", func(t *testing.T) {
		tables := []*Table{
			NewTable([]string{"x", "y"}, 5),
			NewTable([]string{"a", "b", "c"}, 5),
		}
		trainer, err := NewTrainer(&som, tables, &params, rng)
		assert.Nil(t, err)

		trainer.Train(25)

		for _, v := range som.layers[0].data {
			assert.InDelta(t, 0, v, 0.0001)
		}
	})

	t.Run("Train with empty table", func(t *testing.T) {
		tables := []*Table{
			NewTable([]string{"x", "y"}, 5),
			NewTable([]string{"a", "b", "c"}, 5),
		}
		trainer, err := NewTrainer(&som, tables, &params, rng)
		assert.Nil(t, err)

		trainer.Train(10)
	})
}
