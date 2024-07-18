package csv

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/stretchr/testify/assert"
)

func TestSomToCsv(t *testing.T) {
	mockSom := createMockSom()
	var buf bytes.Buffer
	err := SomToCsv(mockSom, &buf, ',', "-")
	assert.NoError(t, err)

	result := buf.String()
	expectedHeader := "node_id,node_x,node_y,CategoricalLayer,a,b,c,d\n"
	fmt.Println(result)
	assert.True(t, strings.HasPrefix(result, expectedHeader))

	lines := strings.Split(result, "\n")
	assert.Equal(t, 5, len(lines)) // Header + 4 nodes

	for i, line := range lines[1:] {
		fields := strings.Split(line, ",")
		assert.Equal(t, 8, len(fields))
		assert.Equal(t, strconv.Itoa(i), fields[0])       // ID
		assert.Contains(t, []string{"0", "1"}, fields[1]) // node_x
		assert.Contains(t, []string{"0", "1"}, fields[2]) // node_y
		assert.Contains(t, []string{"A", "B"}, fields[3]) // CategoricalLayer
		_, err := strconv.ParseFloat(fields[4], 64)       // NumericLayer1
		assert.NoError(t, err)
		_, err = strconv.ParseFloat(fields[5], 64) // NumericLayer2
		assert.NoError(t, err)
	}
}

func createMockSom() *som.Som {
	params := &som.SomConfig{
		Size: layer.Size{Width: 2, Height: 2},
		Layers: []som.LayerDef{
			{
				Name:        "CategoricalLayer",
				Columns:     []string{"A", "B"},
				Metric:      &distance.Hamming{},
				Categorical: true,
				Data:        []float64{0, 1, 0, 1, 1, 0, 1, 0},
			},
			{
				Name:    "NumericLayer1",
				Columns: []string{"a", "b"},
				Metric:  &distance.Euclidean{},
				Data:    []float64{0.1, 0.2, 0.3, 0.4, 1.0, 2.0, 3.0, 4.0},
			},
			{
				Name:    "NumericLayer2",
				Columns: []string{"c", "d"},
				Metric:  &distance.Euclidean{},
				Data:    []float64{1.0, 2.0, 3.0, 4.0, 0.1, 0.2, 0.3, 0.4},
			},
		},
		Neighborhood: &neighborhood.Gaussian{},
	}

	som, err := som.New(params)
	if err != nil {
		panic(err)
	}
	return som
}
