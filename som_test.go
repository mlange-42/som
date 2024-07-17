package som

import (
	"math"
	"testing"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/neighborhood"
	"github.com/stretchr/testify/assert"
)

func TestNewSom(t *testing.T) {
	params := SomConfig{
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
	som := New(&params)

	assert.Equal(t, 2, len(som.layers))
	assert.Equal(t, []float64{0.5, 1.0}, som.weight)
	assert.Equal(t, Size{2, 3}, som.layers[0].size)
}

func TestGetBMU(t *testing.T) {
	params := SomConfig{
		Size: Size{2, 2},
		Layers: []LayerDef{
			{
				Columns: []string{"x", "y"},
				Weight:  0.5,
				Metric:  &distance.Euclidean{},
			},
			{
				Columns: []string{"a", "b"},
				Weight:  1.0,
				Metric:  &distance.SumOfSquares{},
			},
		},
	}
	som := New(&params)

	t.Run("Normal case", func(t *testing.T) {
		data := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		index, dist := som.getBMU(data)
		assert.GreaterOrEqual(t, index, 0)
		assert.Less(t, index, 4)
		assert.GreaterOrEqual(t, dist, 0.0)
	})

	t.Run("Single layer", func(t *testing.T) {
		singleLayerParams := SomConfig{
			Size: Size{1, 1},
			Layers: []LayerDef{
				{
					Columns: []string{"x"},
					Weight:  1.0,
				},
			},
		}
		singleLayerSom := New(&singleLayerParams)
		data := [][]float64{{1.0}}
		index, dist := singleLayerSom.getBMU(data)
		assert.Equal(t, 0, index)
		assert.GreaterOrEqual(t, dist, 0.0)
	})

	t.Run("Large SOM", func(t *testing.T) {
		largeParams := SomConfig{
			Size: Size{10, 10},
			Layers: []LayerDef{
				{
					Columns: []string{"x", "y", "z"},
					Weight:  0.5,
				},
				{
					Columns: []string{"a", "b", "c"},
					Weight:  1.0,
				},
			},
		}
		largeSom := New(&largeParams)
		data := [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}
		index, dist := largeSom.getBMU(data)
		assert.GreaterOrEqual(t, index, 0)
		assert.Less(t, index, 100)
		assert.GreaterOrEqual(t, dist, 0.0)
	})
}

func createSom() *Som {
	params := SomConfig{
		Size: Size{3, 3},
		Layers: []LayerDef{
			{
				Columns: []string{"x", "y"},
				Weight:  0.5,
				Metric:  &distance.Euclidean{},
			},
			{
				Columns: []string{"a", "b"},
				Weight:  1.0,
				Metric:  &distance.SumOfSquares{},
			},
		},
		Neighborhood: &neighborhood.Linear{},
	}
	som := New(&params)

	return &som
}

func TestLearnBasic(t *testing.T) {
	som := createSom()

	t.Run("Basic learning", func(t *testing.T) {
		data := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		initialWeights := make([][][]float64, len(som.layers))
		for l, layer := range som.layers {
			initialWeights[l] = make([][]float64, som.size.Width*som.size.Height)
			for i := range initialWeights[l] {
				initialWeights[l][i] = make([]float64, len(layer.columns))
				copy(initialWeights[l][i], layer.GetNodeAt(i))
			}
		}

		som.learn(data, 0.5, 6.0)

		for l, layer := range som.layers {
			for i := 0; i < som.size.Width*som.size.Height; i++ {
				newWeights := layer.GetNodeAt(i)
				for j := range newWeights {
					assert.False(t, math.IsNaN(newWeights[j]), "Weights should not be NaN")
					assert.NotEqual(t, initialWeights[l][i][j], newWeights[j], "Weights should change after learning")
				}
			}
		}
	})

	t.Run("Zero learning rate", func(t *testing.T) {
		data := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		initialWeights := make([][][]float64, len(som.layers))
		for l, layer := range som.layers {
			initialWeights[l] = make([][]float64, som.size.Width*som.size.Height)
			for i := range initialWeights[l] {
				initialWeights[l][i] = make([]float64, len(layer.columns))
				copy(initialWeights[l][i], layer.GetNodeAt(i))
			}
		}

		som.learn(data, 0.0, 2.0)

		for l, layer := range som.layers {
			for i := 0; i < som.size.Width*som.size.Height; i++ {
				newWeights := layer.GetNodeAt(i)
				for j := range newWeights {
					assert.Equal(t, initialWeights[l][i][j], newWeights[j], "Weights should not change with zero learning rate")
				}
			}
		}
	})
}

func TestLearnRadius(t *testing.T) {
	som := createSom()

	t.Run("Very small radius", func(t *testing.T) {
		data := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		initialWeights := make([][][]float64, len(som.layers))
		for l, layer := range som.layers {
			initialWeights[l] = make([][]float64, som.size.Width*som.size.Height)
			for i := range initialWeights[l] {
				initialWeights[l][i] = make([]float64, len(layer.columns))
				copy(initialWeights[l][i], layer.GetNodeAt(i))
			}
		}

		som.learn(data, 0.5, 0.01)

		for l, layer := range som.layers {
			changedCount := 0
			for i := 0; i < som.size.Width*som.size.Height; i++ {
				newWeights := layer.GetNodeAt(i)
				for j := range newWeights {
					assert.False(t, math.IsNaN(newWeights[j]), "Weights should not be NaN")
					if initialWeights[l][i][j] != newWeights[j] {
						changedCount++
						break
					}
				}
			}
			assert.Equal(t, 1, changedCount, "Only one node (BMU) per layer should change with zero radius")
		}
	})

	t.Run("Very large radius", func(t *testing.T) {
		data := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		som.learn(data, 1.0, 100.0)

		for l, layer := range som.layers {
			lData := data[l]
			for i := 0; i < som.size.Width*som.size.Height; i++ {
				newWeights := layer.GetNodeAt(i)
				for j := range newWeights {
					assert.InDelta(t, lData[j], newWeights[j], 0.5, "All weights should be closer to input data with large radius")
				}
			}
		}
	})
}
