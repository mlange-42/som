package som

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
	"github.com/mlange-42/som/table"
	"github.com/stretchr/testify/assert"
)

func TestPredictorGetBMU(t *testing.T) {
	for _, useTree := range []bool{false, true} {
		t.Run(fmt.Sprintf("UseTree=%t", useTree), func(t *testing.T) {
			conf := SomConfig{
				Size: layer.Size{
					Width:  3,
					Height: 2,
				},
				Neighborhood: &neighborhood.Linear{},
				Layers: []*LayerDef{
					{
						Name:    "L1",
						Columns: []string{"x", "y"},
						Metric:  &distance.Euclidean{},
						Norm:    []norm.Normalizer{&norm.Identity{}, &norm.Identity{}},
						Weights: []float64{
							0, 0, // 0, 0
							0, 1, // 0, 1
							1, 0, // 1, 0
							1, 1, // 1, 1
							2, 0, // 2, 0
							2, 1, // 2, 1
						},
					},
				},
			}

			som, err := New(&conf)
			assert.NoError(t, err)

			tab, err := table.NewWithData([]string{"x", "y"}, []float64{
				-2, 0,
				0, 1,
				1, 0,
				1, 1,
				2, 0,
				2, 2,
			})
			assert.NoError(t, err)

			reader := mockReader{
				Table: tab,
			}
			tables, _, err := conf.PrepareTables(&reader, nil, false, false)
			assert.NoError(t, err)

			p, err := NewPredictor(som, tables, useTree)
			assert.NoError(t, err)

			bmu := p.GetBMUTable()

			assert.Equal(t, 6, bmu.Rows())
			assert.Equal(t, 4, bmu.Columns())
			assert.Equal(t, []float64{
				0, 0, 0, 2,
				1, 0, 1, 0,
				2, 1, 0, 0,
				3, 1, 1, 0,
				4, 2, 0, 0,
				5, 2, 1, 1,
			}, bmu.Data())
		},
		)
	}
}

func BenchmarkPredictorGetBMU_5x5x3_100Rows(b *testing.B) {
	benchmarkPredictorGetBMU(b, 5, 5, 3, 1, false)
}

func BenchmarkPredictorGetBMU_5x5x3_100Rows_Acc(b *testing.B) {
	benchmarkPredictorGetBMU(b, 5, 5, 3, 1, true)
}

func BenchmarkPredictorGetBMU_10x10x5_100Rows(b *testing.B) {
	benchmarkPredictorGetBMU(b, 10, 10, 5, 1, false)
}

func BenchmarkPredictorGetBMU_10x10x5_100Rows_Acc(b *testing.B) {
	benchmarkPredictorGetBMU(b, 10, 10, 5, 1, true)
}

func benchmarkPredictorGetBMU(b *testing.B, width, height, columns, rows int, kdTree bool) {
	b.StopTimer()

	conf := createSomConfig(width, height, columns)

	som, err := New(&conf)
	if err != nil {
		b.Fatal(err)
	}
	som.Randomize(rand.New(rand.NewSource(0)))

	rng := rand.New(rand.NewSource(0))
	tab := table.New([]string{"a", "b", "c", "d", "e"}[:columns], rows)
	for i := range tab.Data() {
		tab.Data()[i] = rng.Float64()*0.2 - 0.1
	}
	var bmu []int
	pred, err := NewPredictor(som, []*table.Table{tab}, kdTree)
	assert.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bmu = pred.GetBMU()
	}
	b.StopTimer()

	assert.Equal(b, rows, len(bmu))
}

func BenchmarkPredictorBMU_5x5x3(b *testing.B) {
	benchmarkPredictorBmu(b, 5, 5, 3, false)
}

func BenchmarkPredictorBMU_5x5x3_Acc(b *testing.B) {
	benchmarkPredictorBmu(b, 5, 5, 3, true)
}

func BenchmarkPredictorBMU_10x10x5(b *testing.B) {
	benchmarkPredictorBmu(b, 10, 10, 5, false)
}

func BenchmarkPredictorBMU_10x10x5_Acc(b *testing.B) {
	benchmarkPredictorBmu(b, 10, 10, 5, true)
}

func benchmarkPredictorBmu(b *testing.B, width, height, columns int, kdTree bool) {
	b.StopTimer()

	conf := createSomConfig(width, height, columns)

	som, err := New(&conf)
	if err != nil {
		b.Fatal(err)
	}
	som.Randomize(rand.New(rand.NewSource(0)))

	tab := table.New([]string{"a", "b", "c", "d", "e"}[:columns], 1)

	pred, err := NewPredictor(som, []*table.Table{tab}, kdTree)
	assert.NoError(b, err)

	data := [][]float64{[]float64{1, 2, 3, 4, 5}[:columns]}
	var bmu int
	var dist float64

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bmu, dist = pred.bmu(data)
	}
	b.StopTimer()

	assert.Less(b, bmu, width*height)
	assert.Less(b, dist, 25.0)
}

func createSomConfig(width, height, columns int) SomConfig {
	return SomConfig{
		Size: layer.Size{
			Width:  width,
			Height: height,
		},
		Neighborhood: &neighborhood.Linear{},
		Layers: []*LayerDef{
			{
				Name:    "L1",
				Columns: []string{"a", "b", "c", "d", "e"}[:columns],
				Metric:  &distance.Euclidean{},
				Norm: []norm.Normalizer{
					&norm.Identity{}, &norm.Identity{}, &norm.Identity{}, &norm.Identity{}, &norm.Identity{},
				}[:columns],
			},
		},
	}
}
