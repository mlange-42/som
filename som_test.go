package som

import (
	"testing"

	"github.com/mlange-42/som/distance"
	"github.com/stretchr/testify/assert"
)

func TestNewSom(t *testing.T) {
	params := SomParams{
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
	assert.Equal(t, []int{0, 2}, som.offset)
	assert.Equal(t, []float64{0.5, 1.0}, som.weight)
	assert.Equal(t, Size{2, 3}, som.layers[0].size)
}

func TestGetBMU(t *testing.T) {
	params := SomParams{
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
		data := []float64{1.0, 2.0, 3.0, 4.0}
		index, dist := som.getBMU(data)
		assert.GreaterOrEqual(t, index, 0)
		assert.Less(t, index, 4)
		assert.GreaterOrEqual(t, dist, 0.0)
	})

	t.Run("Single layer", func(t *testing.T) {
		singleLayerParams := SomParams{
			Size: Size{1, 1},
			Layers: []LayerDef{
				{
					Columns: []string{"x"},
					Weight:  1.0,
				},
			},
		}
		singleLayerSom := New(&singleLayerParams)
		data := []float64{1.0}
		index, dist := singleLayerSom.getBMU(data)
		assert.Equal(t, 0, index)
		assert.GreaterOrEqual(t, dist, 0.0)
	})

	t.Run("Large SOM", func(t *testing.T) {
		largeParams := SomParams{
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
		data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
		index, dist := largeSom.getBMU(data)
		assert.GreaterOrEqual(t, index, 0)
		assert.Less(t, index, 100)
		assert.GreaterOrEqual(t, dist, 0.0)
	})
}
