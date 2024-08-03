package som

import (
	"fmt"
	"testing"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/spatial/kdtree"
)

func TestKDTree(t *testing.T) {
	w, h := 6, 5
	weights := make([]float64, 0, w*h*2)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			weights = append(weights, float64(x), float64(y))
		}
	}

	conf := SomConfig{
		Size: layer.Size{
			Width:  w,
			Height: h,
		},
		Neighborhood: &neighborhood.Linear{},
		Layers: []*LayerDef{
			{
				Name:    "L1",
				Columns: []string{"a", "b"},
				Metric:  &distance.Manhattan{},
				Norm:    []norm.Normalizer{&norm.Identity{}, &norm.Identity{}},
				Weights: weights,
			},
		},
	}

	som, err := New(&conf)
	assert.NoError(t, err)

	locs := newNodeLocations(som)
	tree := kdtree.New(locs, false)

	assert.Equal(t, 30, tree.Count)

	tests := []struct {
		name         string
		data         [][]float64
		expectedIdx  int
		expectedDist float64
	}{}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			tests = append(tests, struct {
				name         string
				data         [][]float64
				expectedIdx  int
				expectedDist float64
			}{
				name:         fmt.Sprintf("Point %d,%d", x, y),
				data:         [][]float64{{float64(x) + 0.1, float64(y) - 0.1}},
				expectedIdx:  som.Size().Index(x, y),
				expectedDist: 0.2,
			})
		}
	}

	for _, test := range tests {
		p := newDataLocation(som, test.data)
		nearest, dist := tree.Nearest(&p)
		bmu := nearest.(*nodeLocation)

		assert.Equal(t, test.expectedIdx, bmu.NodeIndex)
		assert.InDelta(t, test.expectedDist, dist, 1e-10)
	}
}

func TestKDTreeTwoLayers(t *testing.T) {
	conf := SomConfig{
		Size: layer.Size{
			Width:  3,
			Height: 2,
		},
		Neighborhood: &neighborhood.Linear{},
		Layers: []*LayerDef{
			{
				Name:    "L1",
				Columns: []string{"a", "b"},
				Metric:  &distance.Euclidean{},
				Norm:    []norm.Normalizer{&norm.Identity{}, &norm.Identity{}},
				Weights: []float64{
					0, 0,
					0, 1,
					1, 0,
					1, 1,
					2, 0,
					2, 1,
				},
			},
			{
				Name:    "L2",
				Columns: []string{"x", "y"},
				Metric:  &distance.Euclidean{},
				Norm:    []norm.Normalizer{&norm.Identity{}, &norm.Identity{}},
				Weights: []float64{
					0, 0,
					0, 1,
					1, 0,
					1, 1,
					2, 0,
					2, 1,
				},
			},
		},
	}

	som, err := New(&conf)
	assert.NoError(t, err)

	locs := newNodeLocations(som)
	tree := kdtree.New(locs, false)

	assert.Equal(t, 6, tree.Count)

	tests := []struct {
		name         string
		data         [][]float64
		expectedIdx  int
		expectedDist float64
	}{
		{
			name:         "Data point at (0,0), (0,0)",
			data:         [][]float64{{0, 0}, {0, 0}},
			expectedIdx:  0,
			expectedDist: 0.0,
		},
		{
			name:         "Data point at (1,1), (1,1)",
			data:         [][]float64{{1, 1}, {1, 1}},
			expectedIdx:  3,
			expectedDist: 0.0,
		},
		{
			name:         "Data point at (1,2), (1,2)",
			data:         [][]float64{{1, 2}, {1, 2}},
			expectedIdx:  3,
			expectedDist: 2.0,
		},
		{
			name:         "Data point at (3,1), (3,1)",
			data:         [][]float64{{4, 1}, {4, 1}},
			expectedIdx:  5,
			expectedDist: 4.0,
		},
	}

	for _, test := range tests {
		p := newDataLocation(som, test.data)
		nearest, dist := tree.Nearest(&p)
		bmu := nearest.(*nodeLocation)

		assert.Equal(t, test.expectedIdx, bmu.NodeIndex)
		assert.InDelta(t, test.expectedDist, dist, 1e-10)
	}
}
