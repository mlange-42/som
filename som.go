package som

import (
	"math"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/neighborhood"
)

type SomParams struct {
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
	size   Size
	layers []Layer
	offset []int
	weight []float64
	metric []distance.Distance
}

func New(params *SomParams) Som {
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
		size:   params.Size,
		layers: lay,
		offset: offset,
		weight: weight,
		metric: metric,
	}
}

func (s *Som) learn(data []float64, alpha, radius float64) {

}

func (s *Som) getBMU(data []float64) (int, float64) {
	units := s.size.Width * s.size.Height

	minDist := math.MaxFloat64
	minIndex := -1
	for i := 0; i < units; i++ {
		totalDist := 0.0
		for l, layer := range s.layers {
			node := layer.GetNodeAt(l)
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
