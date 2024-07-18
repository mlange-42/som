package som

import (
	"testing"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/table"
	"github.com/stretchr/testify/assert"
)

func TestPredictorGetBMU(t *testing.T) {
	conf := SomConfig{
		Size: layer.Size{
			Width:  3,
			Height: 2,
		},
		Neighborhood: &neighborhood.Linear{},
		Layers: []LayerDef{
			{
				Name:    "L1",
				Columns: []string{"x", "y"},
				Metric:  &distance.Euclidean{},
				Data: []float64{
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
		0, 0,
		0, 1,
		1, 0,
		1, 1,
		2, 0,
		2, 1,
	})
	assert.NoError(t, err)

	reader := mockReader{
		Table: tab,
	}
	tables, err := conf.PrepareTables(&reader, false)
	assert.NoError(t, err)

	p, err := NewPredictor(&som, tables)
	assert.NoError(t, err)

	bmu, err := p.GetBMU()
	assert.NoError(t, err)

	assert.Equal(t, 6, bmu.Rows())
	assert.Equal(t, 3, bmu.Columns())
	assert.Equal(t, []float64{
		0, 0, 0,
		1, 0, 1,
		2, 1, 0,
		3, 1, 1,
		4, 2, 0,
		5, 2, 1,
	}, bmu.Data())
}
