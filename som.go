package som

import (
	"fmt"
	"math"

	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
	"github.com/mlange-42/som/table"
)

type SomConfig struct {
	Size         layer.Size
	Layers       []LayerDef
	Neighborhood neighborhood.Neighborhood
}

// PrepareTables reads the CSV data and creates a table for each layer defined in the SomConfig.
// If a categorical layer has no columns specified, it will attempt to read the class names for that layer
// and create a table from the classes. The created tables are returned in the same order as
// the layers in the SomConfig.
func (c *SomConfig) PrepareTables(reader csv.Reader) ([]*table.Table, error) {
	tables := make([]*table.Table, len(c.Layers))
	for i := range c.Layers {
		layer := &c.Layers[i]
		if len(layer.Columns) == 0 {
			if !layer.Categorical {
				return nil, fmt.Errorf("layer %d has no columns", i)
			}
			classes, err := reader.ReadLabels(layer.Name)
			if err != nil {
				return nil, err
			}
			table := conv.ClassesToTable(classes)
			layer.Columns = table.ColumnNames()
			tables[i] = table
			continue
		}
		table, err := reader.ReadColumns(layer.Columns)
		if err != nil {
			return nil, err
		}
		tables[i] = table
	}
	return tables, nil
}

type LayerDef struct {
	Name        string
	Columns     []string
	Norm        []norm.Normalizer
	Metric      distance.Distance
	Weight      float64
	Categorical bool
	Data        []float64
}

type Som struct {
	size         layer.Size
	layers       []layer.Layer
	neighborhood neighborhood.Neighborhood
}

func New(params *SomConfig) (Som, error) {
	lay := make([]layer.Layer, len(params.Layers))
	for i, l := range params.Layers {
		if len(l.Columns) == 0 {
			return Som{}, fmt.Errorf("layer %d has no columns", i)
		}
		weight := l.Weight
		if weight == 0 {
			weight = 1
		}
		metric := l.Metric
		if metric == nil {
			if l.Categorical {
				metric = &distance.Hamming{}
			} else {
				metric = &distance.Euclidean{}
			}
		}
		if len(l.Data) == 0 {
			lay[i] = layer.New(l.Name, l.Columns, l.Norm, params.Size, metric, weight, l.Categorical)
		} else {
			var err error
			lay[i], err = layer.NewWithData(l.Name, l.Columns, l.Norm, params.Size, metric, weight, l.Categorical, l.Data)
			if err != nil {
				return Som{}, err
			}
		}
	}
	return Som{
		size:         params.Size,
		layers:       lay,
		neighborhood: params.Neighborhood,
	}, nil
}

func (s *Som) Size() layer.Size {
	return s.size
}

func (s *Som) Neighborhood() neighborhood.Neighborhood {
	return s.neighborhood
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

	for l, lay := range s.layers {
		lData := data[l]
		cols := lay.Columns()

		for x := xMin; x <= xMax; x++ {
			for y := yMin; y <= yMax; y++ {
				node := lay.GetNode(x, y)
				r := s.neighborhood.Weight(x, y, xBmu, yBmu, radius)
				for i := 0; i < cols; i++ {
					if math.IsNaN(lData[i]) {
						continue
					}
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
			dist := layer.Metric().Distance(node, data[l])
			totalDist += layer.Weight() * dist
		}
		if totalDist < minDist {
			minDist = totalDist
			minIndex = i
		}
	}

	return minIndex, minDist
}

func (s *Som) Layers() []layer.Layer {
	return s.layers
}
