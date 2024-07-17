package som

import (
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
	Columns []string
	Metric  distance.Distance
	Weight  float64
}

type Som struct {
	size         Size
	layers       []Layer
	offset       []int
	weight       []float64
	metric       []distance.Distance
	neighborhood neighborhood.Neighborhood
}

func New(params *SomConfig) Som {
	lay := make([]Layer, len(params.Layers))
	weight := make([]float64, len(params.Layers))
	offset := make([]int, len(params.Layers))
	metric := make([]distance.Distance, len(params.Layers))
	off := 0
	for i, l := range params.Layers {
		lay[i] = NewLayer(l.Columns, params.Size)
		offset[i] = off
		off += len(l.Columns)

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
		offset:       offset,
		weight:       weight,
		metric:       metric,
		neighborhood: params.Neighborhood,
	}
}

func (s *Som) learn(data []float64, alpha, radius float64) {
	lim := s.neighborhood.MaxRadius(radius)
	if lim < 0 {
		lim = s.size.Width * s.size.Height
	}

	bmuIdx, _ := s.getBMU(data)
	xBmu, yBmu := s.size.CoordsAt(bmuIdx)
	xMin, yMin := max(xBmu-lim, 0), max(yBmu-lim, 0)
	xMax, yMax := min(xBmu+lim, s.size.Width-1), min(yBmu+lim, s.size.Height-1)

	for l, layer := range s.layers {
		offset := s.offset[l]
		cols := len(layer.columns)

		for x := xMin; x <= xMax; x++ {
			for y := yMin; y <= yMax; y++ {
				node := layer.GetNode(x, y)
				r := s.neighborhood.Weight(x, y, xBmu, yBmu, radius)
				for i := 0; i < cols; i++ {
					node[i] += alpha * r * (data[offset+i] - node[i])
				}
			}
		}
	}
}

func (s *Som) getBMU(data []float64) (int, float64) {
	units := s.size.Width * s.size.Height

	minDist := math.MaxFloat64
	minIndex := -1
	for i := 0; i < units; i++ {
		totalDist := 0.0
		for l, layer := range s.layers {
			node := layer.GetNodeAt(i)
			offset := s.offset[l]
			cols := len(layer.columns)
			dist := s.metric[l].Distance(node, data[offset:offset+cols])
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
