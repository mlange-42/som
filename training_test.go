package som

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/table"
	"github.com/stretchr/testify/assert"
)

func TestNewTrainer(t *testing.T) {
	params := TrainingConfig{}
	somParams := SomConfig{
		Size: layer.Size{Width: 2, Height: 3},
		Layers: []*LayerDef{
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

	t1 := []*table.Table{
		table.New([]string{"x", "y"}, 5),
		table.New([]string{"a", "b", "c"}, 5),
	}
	_, err = NewTrainer(som, t1, &params, rng)
	assert.Nil(t, err)

	t1 = []*table.Table{
		table.New([]string{"x", "y"}, 5),
		table.New([]string{"a", "b"}, 5),
	}
	_, err = NewTrainer(som, t1, &params, rng)
	assert.NotNil(t, err)

	t1 = []*table.Table{
		table.New([]string{"x", "y"}, 5),
		table.New([]string{"a", "b", "c", "d"}, 5),
	}
	_, err = NewTrainer(som, t1, &params, rng)
	assert.NotNil(t, err)

	t1 = []*table.Table{
		table.New([]string{"x", "y"}, 5),
		table.New([]string{"a", "c", "b"}, 5),
	}
	_, err = NewTrainer(som, t1, &params, rng)
	assert.NotNil(t, err)
}

func TestTrainerDecay(t *testing.T) {
	params := TrainingConfig{
		Epochs:             100,
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Power{Start: 5, End: 0.5},
	}
	somParams := SomConfig{
		Size: layer.Size{Width: 2, Height: 3},
		Layers: []*LayerDef{
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
		MapMetric:    &neighborhood.Manhattan{},
	}
	som, err := New(&somParams)
	assert.NoError(t, err)

	rng := rand.New(rand.NewSource(1))

	t1 := []*table.Table{
		table.New([]string{"x", "y"}, 5),
		table.New([]string{"a", "b", "c"}, 5),
	}
	trainer, err := NewTrainer(som, t1, &params, rng)
	assert.Nil(t, err)

	progress := make(chan TrainingProgress)

	go func() {
		trainer.Train(progress)
	}()

	for epoch := range progress {
		fmt.Println(epoch)
	}
}

func TestTrainerTrain(t *testing.T) {
	params := TrainingConfig{
		Epochs:             0,
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Power{Start: 3, End: 0.5},
	}
	somParams := SomConfig{
		Size: layer.Size{Width: 2, Height: 3},
		Layers: []*LayerDef{
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
		MapMetric:    &neighborhood.Manhattan{},
	}
	som, err := New(&somParams)
	assert.NoError(t, err)

	rng := rand.New(rand.NewSource(1))

	t.Run("Train with zero epochs", func(t *testing.T) {
		tables := []*table.Table{
			table.New([]string{"x", "y"}, 5),
			table.New([]string{"a", "b", "c"}, 5),
		}
		trainer, err := NewTrainer(som, tables, &params, rng)
		assert.Nil(t, err)

		progress := make(chan TrainingProgress)

		go func() {
			trainer.Train(progress)
		}()

		for range progress {
		}

		for _, v := range som.layers[0].Data() {
			assert.NotEqual(t, 0, v)
		}
	})

	t.Run("Train with one epoch", func(t *testing.T) {
		tables := []*table.Table{
			table.New([]string{"x", "y"}, 5),
			table.New([]string{"a", "b", "c"}, 5),
		}
		p := params
		p.Epochs = 1
		trainer, err := NewTrainer(som, tables, &p, rng)
		assert.Nil(t, err)

		progress := make(chan TrainingProgress)

		go func() {
			trainer.Train(progress)
		}()

		for range progress {
		}
	})

	t.Run("Train with multiple epochs", func(t *testing.T) {
		tables := []*table.Table{
			table.New([]string{"x", "y"}, 5),
			table.New([]string{"a", "b", "c"}, 5),
		}
		p := params
		p.Epochs = 25
		trainer, err := NewTrainer(som, tables, &p, rng)
		assert.Nil(t, err)

		progress := make(chan TrainingProgress)

		go func() {
			trainer.Train(progress)
		}()

		for range progress {
		}

		for _, v := range som.layers[0].Data() {
			assert.InDelta(t, 0, v, 0.0001)
		}
	})

	t.Run("Train with empty table", func(t *testing.T) {
		tables := []*table.Table{
			table.New([]string{"x", "y"}, 5),
			table.New([]string{"a", "b", "c"}, 5),
		}
		p := params
		p.Epochs = 25
		trainer, err := NewTrainer(som, tables, &p, rng)
		assert.Nil(t, err)

		progress := make(chan TrainingProgress)

		go func() {
			trainer.Train(progress)
		}()

		for range progress {
		}
	})
}
