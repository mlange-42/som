package som

import (
	"fmt"
	"math"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/neighborhood"
)

type SomConfig struct {
	Size         Size
	Layers       []LayerDef
	Neighborhood neighborhood.Neighborhood
}

type LayerDef struct {
	Name        string
	Columns     []string
	Metric      distance.Distance
	Weight      float64
	Categorical bool
}

type Som struct {
	size         Size
	layers       []Layer
	weight       []float64
	metric       []distance.Distance
	neighborhood neighborhood.Neighborhood
}

func New(params *SomConfig) (Som, error) {
	lay := make([]Layer, len(params.Layers))
	weight := make([]float64, len(params.Layers))
	metric := make([]distance.Distance, len(params.Layers))
	for i, l := range params.Layers {
		if len(l.Columns) == 0 {
			return Som{}, fmt.Errorf("layer %d has no columns", i)
		}

		lay[i] = NewLayer(l.Name, l.Columns, params.Size, l.Categorical)

		weight[i] = l.Weight
		if weight[i] == 0 {
			weight[i] = 1
		}

		metric[i] = l.Metric
		if metric[i] == nil {
			metric[i] = &distance.Euclidean{}
		}
	}
	return Som{
		size:         params.Size,
		layers:       lay,
		weight:       weight,
		metric:       metric,
		neighborhood: params.Neighborhood,
	}, nil
}

func (s *Som) learn(data [][]float64, alpha, radius float64) {
	lim := s.neighborhood.MaxRadius(radius)
	if lim < 0 {
		lim = s.size.Width * s.size.Height
	}

	bmuIdx, _ := s.getBMU(data)
	xBmu, yBmu := s.size.CoordsAt(bmuIdx)
	xMin, yMin := max(xBmu-lim, 0), max(yBmu-lim, 0)
	xMax, yMax := min(xBmu+lim, s.size.Width-1), min(yBmu+lim, s.size.Height-1)

	for l, layer := range s.layers {
		lData := data[l]
		cols := len(layer.columns)

		for x := xMin; x <= xMax; x++ {
			for y := yMin; y <= yMax; y++ {
				node := layer.GetNode(x, y)
				r := s.neighborhood.Weight(x, y, xBmu, yBmu, radius)
				for i := 0; i < cols; i++ {
					node[i] += alpha * r * (lData[i] - node[i])
				}
			}
		}
	}
}

func (s *Som) getBMU(data [][]float64) (int, float64) {
	units := s.size.Width * s.size.Height

	minDist := math.MaxFloat64
	minIndex := -1
	for i := 0; i < units; i++ {
		totalDist := 0.0
		for l, layer := range s.layers {
			node := layer.GetNodeAt(i)
			dist := s.metric[l].Distance(node, data[l])
			totalDist += s.weight[l] * dist
		}
		if totalDist < minDist {
			minDist = totalDist
			minIndex = i
		}
	}

	return minIndex, minDist
}

func (s *Som) Layers() []Layer {
	return s.layers
}
