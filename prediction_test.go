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

func BenchmarkPredictorGetBMU_20x30x5_1kRows(b *testing.B) {
	benchmarkPredictorGetBMU(b, 20, 30, 1000, false)
}

func BenchmarkPredictorGetBMU_20x30x5_1kRows_Acc(b *testing.B) {
	benchmarkPredictorGetBMU(b, 20, 30, 1000, true)
}

func BenchmarkPredictorGetBMU_100x100x5_1kRows(b *testing.B) {
	benchmarkPredictorGetBMU(b, 100, 100, 1000, false)
}

func BenchmarkPredictorGetBMU_100x100x5_1kRows_Acc(b *testing.B) {
	benchmarkPredictorGetBMU(b, 100, 100, 1000, true)
}

func benchmarkPredictorGetBMU(b *testing.B, width, height, rows int, kdTree bool) {
	b.StopTimer()

	conf := SomConfig{
		Size: layer.Size{
			Width:  width,
			Height: height,
		},
		Neighborhood: &neighborhood.Linear{},
		Layers: []*LayerDef{
			{
				Name:    "L1",
				Columns: []string{"a", "b", "c", "d", "e"},
				Metric:  &distance.Euclidean{},
				Norm:    []norm.Normalizer{&norm.Identity{}, &norm.Identity{}, &norm.Identity{}, &norm.Identity{}, &norm.Identity{}},
			},
		},
	}

	som, err := New(&conf)
	if err != nil {
		b.Fatal(err)
	}
	som.Randomize(rand.New(rand.NewSource(0)))

	rng := rand.New(rand.NewSource(0))
	tab := table.New([]string{"a", "b", "c", "d", "e"}, rows)
	for i := range tab.Data() {
		tab.Data()[i] = rng.Float64()*2 - 1
	}
	var bmu []int
	pred, err := NewPredictor(som, []*table.Table{tab}, kdTree)
	assert.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bmu = pred.GetBMU()
	}
	b.StopTimer()
	_ = bmu
}
