package som

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/table"
	"github.com/stretchr/testify/assert"
)

type mockReader struct {
	Table  *table.Table
	Labels []string
}

func (r *mockReader) ReadColumns(columns []string) (*table.Table, error) {
	return r.Table, nil
}

func (r *mockReader) ReadLabels(column string) ([]string, error) {
	return r.Labels, nil
}

func (r *mockReader) NoData() string {
	return ""
}

func TestNew(t *testing.T) {
	t.Run("Valid configuration", func(t *testing.T) {
		params := &SomConfig{
			Size: layer.Size{Width: 3, Height: 3},
			Layers: []*LayerDef{
				{
					Name:    "Layer1",
					Columns: []string{"x", "y"},
					Weight:  0.5,
					Metric:  &distance.Manhattan{},
				},
				{
					Name:    "Layer2",
					Columns: []string{"a", "b", "c"},
					Weight:  1.0,
					Metric:  &distance.Euclidean{},
				},
			},
			Neighborhood: &neighborhood.Gaussian{},
		}

		som, err := New(params)
		assert.NoError(t, err)
		assert.Equal(t, params.Size, som.size)
		assert.Len(t, som.layers, 2)
		assert.IsType(t, &neighborhood.Gaussian{}, som.neighborhood)

		assert.Equal(t, "Layer1", som.layers[0].Name())
		assert.Equal(t, "Layer2", som.layers[1].Name())
		assert.Equal(t, 0.5, som.layers[0].Weight())
		assert.Equal(t, 1.0, som.layers[1].Weight())
		assert.IsType(t, &distance.Manhattan{}, som.layers[0].Metric())
		assert.IsType(t, &distance.Euclidean{}, som.layers[1].Metric())
	})

	t.Run("Categorical with reader", func(t *testing.T) {
		params := &SomConfig{
			Size: layer.Size{Width: 3, Height: 3},
			Layers: []*LayerDef{
				{
					Name:    "Layer1",
					Columns: []string{"x", "y"},
					Weight:  0.5,
					Metric:  &distance.Manhattan{},
				},
				{
					Name:        "Layer2",
					Columns:     nil,
					Weight:      1.0,
					Metric:      &distance.Hamming{},
					Categorical: true,
				},
			},
			Neighborhood: &neighborhood.Gaussian{},
		}

		tab, err := table.NewWithData([]string{"x", "y"}, []float64{1, 2, 4, 5, 7, 8})
		assert.NoError(t, err)
		reader := mockReader{
			Table:  tab,
			Labels: []string{"A", "B", "A"},
		}
		tables, _, err := params.PrepareTables(&reader, nil, false, false)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(tables))
		assert.Equal(t, []string{"x", "y"}, tables[0].ColumnNames())
		assert.Equal(t, []string{"A", "B"}, tables[1].ColumnNames())

		assert.Equal(t, []float64{1, 2, 4, 5, 7, 8}, tables[0].Data())
		assert.Equal(t, []float64{1, 0, 0, 1, 1, 0}, tables[1].Data())

		som, err := New(params)
		assert.NoError(t, err)
		assert.Equal(t, params.Size, som.size)
		assert.Len(t, som.layers, 2)
		assert.IsType(t, &neighborhood.Gaussian{}, som.neighborhood)

		assert.Equal(t, []string{"A", "B"}, som.layers[1].ColumnNames())
	})

	t.Run("Empty columns", func(t *testing.T) {
		params := &SomConfig{
			Size: layer.Size{Width: 2, Height: 2},
			Layers: []*LayerDef{
				{
					Name:        "EmptyLayer",
					Categorical: true,
				},
			},
		}

		_, err := New(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "layer EmptyLayer has no columns")
	})

	t.Run("Default weight and metric", func(t *testing.T) {
		params := &SomConfig{
			Size: layer.Size{Width: 2, Height: 2},
			Layers: []*LayerDef{
				{
					Name:    "DefaultLayer",
					Columns: []string{"x"},
				},
			},
		}

		som, err := New(params)
		assert.NoError(t, err)
		_ = som
	})

	t.Run("Multiple layers with different configurations", func(t *testing.T) {
		params := &SomConfig{
			Size: layer.Size{Width: 4, Height: 4},
			Layers: []*LayerDef{
				{
					Name:        "CategoricalLayer",
					Categorical: true,
					Columns:     []string{"category1", "category2"},
					Weight:      0.7,
				},
				{
					Name:    "NumericLayer",
					Columns: []string{"x", "y", "z"},
					Weight:  1.2,
					Metric:  &distance.Manhattan{},
				},
			},
			Neighborhood: &neighborhood.Linear{},
		}

		som, err := New(params)
		assert.NoError(t, err)
		assert.Len(t, som.layers, 2)
		assert.True(t, som.layers[0].IsCategorical())
		assert.False(t, som.layers[1].IsCategorical())
		assert.IsType(t, &neighborhood.Linear{}, som.neighborhood)
	})

	t.Run("Invalid layer configuration", func(t *testing.T) {
		params := &SomConfig{
			Size: layer.Size{Width: 2, Height: 2},
			Layers: []*LayerDef{
				{
					Name: "InvalidLayer",
				},
			},
		}

		_, err := New(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "layer InvalidLayer has no columns")
	})
}

func TestGetBMU(t *testing.T) {
	params := SomConfig{
		Size: layer.Size{Width: 2, Height: 2},
		Layers: []*LayerDef{
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
	som, err := New(&params)
	assert.NoError(t, err)

	t.Run("Normal case", func(t *testing.T) {
		data := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		index, dist := som.GetBMU(data)
		assert.GreaterOrEqual(t, index, 0)
		assert.Less(t, index, 4)
		assert.GreaterOrEqual(t, dist, 0.0)
	})

	t.Run("Single layer", func(t *testing.T) {
		singleLayerParams := SomConfig{
			Size: layer.Size{Width: 1, Height: 1},
			Layers: []*LayerDef{
				{
					Columns: []string{"x"},
					Weight:  1.0,
				},
			},
		}
		singleLayerSom, err := New(&singleLayerParams)
		assert.NoError(t, err)

		data := [][]float64{{1.0}}
		index, dist := singleLayerSom.GetBMU(data)
		assert.Equal(t, 0, index)
		assert.GreaterOrEqual(t, dist, 0.0)
	})

	t.Run("Large SOM", func(t *testing.T) {
		largeParams := SomConfig{
			Size: layer.Size{Width: 10, Height: 10},
			Layers: []*LayerDef{
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
		largeSom, err := New(&largeParams)
		assert.NoError(t, err)

		data := [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}
		index, dist := largeSom.GetBMU(data)
		assert.GreaterOrEqual(t, index, 0)
		assert.Less(t, index, 100)
		assert.GreaterOrEqual(t, dist, 0.0)
	})
}

func createSom() *Som {
	params := SomConfig{
		Size: layer.Size{Width: 3, Height: 3},
		Layers: []*LayerDef{
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
		MapMetric:    &neighborhood.ManhattanMetric{},
	}

	som, err := New(&params)
	if err != nil {
		panic(err)
	}

	return som
}

func TestLearnBasic(t *testing.T) {
	som := createSom()

	t.Run("Basic learning", func(t *testing.T) {
		data := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
		initialWeights := make([][][]float64, len(som.layers))
		for l, layer := range som.layers {
			initialWeights[l] = make([][]float64, som.size.Width*som.size.Height)
			for i := range initialWeights[l] {
				initialWeights[l][i] = make([]float64, len(layer.ColumnNames()))
				copy(initialWeights[l][i], layer.GetNodeAt(i))
			}
		}

		som.Learn(data, 0.5, 6.0, 0.0)

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
				initialWeights[l][i] = make([]float64, len(layer.ColumnNames()))
				copy(initialWeights[l][i], layer.GetNodeAt(i))
			}
		}

		som.Learn(data, 0.0, 2.0, 0.0)

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
				initialWeights[l][i] = make([]float64, len(layer.ColumnNames()))
				copy(initialWeights[l][i], layer.GetNodeAt(i))
			}
		}

		som.Learn(data, 0.5, 0.01, 0.0)

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
		som.Learn(data, 1.0, 100.0, 0.0)

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

func TestNodeDistance(t *testing.T) {
	config := SomConfig{
		Size: layer.Size{Width: 2, Height: 2},
		Layers: []*LayerDef{
			{
				Columns: []string{"x", "y"},
				Weight:  1.0,
				Metric:  &distance.Euclidean{},
				Weights: []float64{
					0.0, 0.0,
					0.0, 2.0,
					2.0, 0.0,
					2.0, 2.0,
				},
			},
		},
	}
	som, err := New(&config)
	if err != nil {
		panic(err)
	}

	assert.InDelta(t, 0, som.nodeDistance(0, 0), 0.001)
	assert.InDelta(t, 2, som.nodeDistance(0, 1), 0.001)
	assert.InDelta(t, 2, som.nodeDistance(0, 2), 0.001)
	assert.InDelta(t, math.Sqrt(8), som.nodeDistance(0, 3), 0.001)

	assert.InDelta(t, math.Sqrt(8), som.nodeDistance(1, 2), 0.001)
}

func TestNodeMapDistance(t *testing.T) {
	config := SomConfig{
		Size:      layer.Size{Width: 2, Height: 2},
		MapMetric: &neighborhood.EuclideanMetric{},
		Layers: []*LayerDef{
			{
				Columns: []string{"x", "y"},
				Weight:  1.0,
				Metric:  &distance.Euclidean{},
				Weights: []float64{
					0.0, 0.0,
					0.0, 2.0,
					2.0, 0.0,
					2.0, 2.0,
				},
			},
		},
	}
	som, err := New(&config)
	if err != nil {
		panic(err)
	}

	assert.InDelta(t, 1, som.MapMetric().Distance(0, 0, 1, 0), 0.001)
	assert.InDelta(t, 1, som.MapMetric().Distance(0, 1, 1, 1), 0.001)
	assert.InDelta(t, math.Sqrt(2), som.MapMetric().Distance(0, 0, 1, 1), 0.001)
}

func BenchmarkGetBMU_5x5x3(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(5, 5, 3, &neighborhood.Gaussian{})
	data := [][]float64{{1.0, 2.0, 3.0}}

	var bmu int

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bmu, _ = som.GetBMU(data)
	}
	b.StopTimer()

	assert.Less(b, bmu, 25)
}

func BenchmarkGetBMU_10x10x5(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(10, 10, 5, &neighborhood.Gaussian{})
	data := [][]float64{{1.0, 2.0, 3.0, 4.0, 5.0}}

	var bmu int

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bmu, _ = som.GetBMU(data)
	}
	b.StopTimer()

	assert.Less(b, bmu, 100)
}

func BenchmarkUpdateWeights_5x5x3_Gaussian2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(5, 5, 3, &neighborhood.Gaussian{})
	data := [][]float64{{1.0, 2.0, 3.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.0)
	}
}

func BenchmarkUpdateWeights_10x10x5_Gaussian2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(10, 10, 5, &neighborhood.Gaussian{})
	data := [][]float64{{1.0, 2.0, 3.0, 4.0, 5.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.0)
	}
}

func BenchmarkUpdateWeights_5x5x3_Linear2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(5, 5, 3, &neighborhood.Linear{})
	data := [][]float64{{1.0, 2.0, 3.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.0)
	}
}

func BenchmarkUpdateWeights_10x10x5_Linear2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(10, 10, 5, &neighborhood.Linear{})
	data := [][]float64{{1.0, 2.0, 3.0, 4.0, 5.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.0)
	}
}

func BenchmarkUpdateWeightsViSOM_5x5x3_Gaussian2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(5, 5, 3, &neighborhood.Gaussian{})
	data := [][]float64{{1.0, 2.0, 3.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.1)
	}
}

func BenchmarkUpdateWeightsViSOM_10x10x5_Gaussian2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(10, 10, 5, &neighborhood.Gaussian{})
	data := [][]float64{{1.0, 2.0, 3.0, 4.0, 5.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.1)
	}
}

func BenchmarkUpdateWeightsViSOM_5x5x3_Linear2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(5, 5, 3, &neighborhood.Linear{})
	data := [][]float64{{1.0, 2.0, 3.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.1)
	}
}

func BenchmarkUpdateWeightsViSOM_10x10x5_Linear2(b *testing.B) {
	b.StopTimer()
	som := createBenchSom(10, 10, 5, &neighborhood.Linear{})
	data := [][]float64{{1.0, 2.0, 3.0, 4.0, 5.0}}

	bmu, _ := som.GetBMU(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		som.updateWeights(bmu, data, 0.01, 2.0, 0.1)
	}
}

func createBenchSom(width, height int, dims int, neigh neighborhood.Neighborhood) *Som {
	cols := make([]string, dims)
	for i := 0; i < dims; i++ {
		cols[i] = fmt.Sprintf("x%d", i)
	}

	params := SomConfig{
		Size: layer.Size{Width: width, Height: height},
		Layers: []*LayerDef{
			{
				Columns: cols,
				Metric:  &distance.Euclidean{},
			},
		},
		Neighborhood: neigh,
		MapMetric:    &neighborhood.ManhattanMetric{},
	}

	som, err := New(&params)
	if err != nil {
		panic(err)
	}

	som.Randomize(rand.New(rand.NewSource(0)))

	return som
}

func TestUMatrix(t *testing.T) {
	params := &SomConfig{
		Size: layer.Size{Width: 2, Height: 2},
		Layers: []*LayerDef{
			{
				Name:    "Layer1",
				Columns: []string{"x", "y"},
				Weight:  1.0,
				Metric:  &distance.Euclidean{},
			},
		},
		Neighborhood: &neighborhood.Gaussian{},
	}

	som, err := New(params)
	assert.NoError(t, err)

	// Set known weights for testing
	som.layers[0].Set(1, 0, 0, 1)
	som.layers[0].Set(0, 1, 1, 2)
	som.layers[0].Set(1, 1, 0, 1)
	som.layers[0].Set(1, 1, 1, 2)

	uMatrix := som.UMatrix(true)

	assert.Equal(t, 3, len(uMatrix))
	assert.Equal(t, 3, len(uMatrix[0]))

	// Check horizontal distances
	assert.InDelta(t, 1.0, uMatrix[0][1], 0.001)
	assert.InDelta(t, 1.0, uMatrix[2][1], 0.001)

	// Check vertical distances
	assert.InDelta(t, 2.0, uMatrix[1][0], 0.001)
	assert.InDelta(t, 2.0, uMatrix[1][2], 0.001)

	// Node center (average of surrounding distances)
	assert.InDelta(t, 1.5, uMatrix[0][0], 0.001)

	// Check center (average of surrounding distances)
	assert.InDelta(t, 1.5, uMatrix[1][1], 0.001)
}
