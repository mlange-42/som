package som

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/mlange-42/som/conv"
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
func (c *SomConfig) PrepareTables(reader table.Reader, updateNormalizers bool) ([]*table.Table, error) {
	tables := make([]*table.Table, len(c.Layers))
	for i := range c.Layers {
		layer := &c.Layers[i]

		if layer.Categorical {
			classes, err := reader.ReadLabels(layer.Name)
			if err != nil {
				return nil, err
			}

			table, err := conv.ClassesToTable(classes, layer.Columns)
			if err != nil {
				return nil, err
			}
			layer.Columns = table.ColumnNames()
			tables[i] = table

			continue
		}

		if len(layer.Columns) == 0 {
			return nil, fmt.Errorf("layer %s has no columns", layer.Name)
		}

		table, err := reader.ReadColumns(layer.Columns)
		if err != nil {
			return nil, err
		}

		if len(layer.Norm) != 0 {
			for j := range layer.Columns {
				if updateNormalizers {
					layer.Norm[j].Initialize(table, j)
				}
				table.NormalizeColumn(j, layer.Norm[j])
			}
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
	layers       []*layer.Layer
	neighborhood neighborhood.Neighborhood
}

func New(params *SomConfig) (*Som, error) {
	lay := make([]*layer.Layer, len(params.Layers))
	for i, l := range params.Layers {
		if len(l.Columns) == 0 {
			return nil, fmt.Errorf("layer %s has no columns", l.Name)
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

		var err error
		if len(l.Data) == 0 {
			lay[i], err = layer.New(l.Name, l.Columns, l.Norm, params.Size, metric, weight, l.Categorical)
		} else {
			lay[i], err = layer.NewWithData(l.Name, l.Columns, l.Norm, params.Size, metric, weight, l.Categorical, l.Data)
		}
		if err != nil {
			return nil, err
		}
	}
	return &Som{
		size:         params.Size,
		layers:       lay,
		neighborhood: params.Neighborhood,
	}, nil
}

func (s *Som) Size() *layer.Size {
	return &s.size
}

func (s *Som) Neighborhood() neighborhood.Neighborhood {
	return s.neighborhood
}

func (s *Som) learn(data [][]float64, alpha, radius, lambda float64) float64 {
	bmuIdx, dist := s.getBMU(data)

	if lambda > 0 {
		s.updateWeightsVI(bmuIdx, data, alpha, radius, lambda)
	} else {
		s.updateWeights(bmuIdx, data, alpha, radius)
	}

	return dist
}

func (s *Som) getBMU(data [][]float64) (int, float64) {
	units := s.size.Width * s.size.Height

	minDist := math.MaxFloat64
	minIndex := -1
	for i := 0; i < units; i++ {
		totalDist := s.distance(data, i)
		if totalDist < minDist {
			minDist = totalDist
			minIndex = i
		}
	}

	return minIndex, minDist
}

func (s *Som) updateWeights(bmuIdx int, data [][]float64, alpha, radius float64) {
	lim := s.neighborhood.MaxRadius(radius)
	if lim < 0 {
		lim = s.size.Width * s.size.Height
	}

	xBmu, yBmu := s.size.CoordsAt(bmuIdx)
	xMin, yMin := max(xBmu-lim, 0), max(yBmu-lim, 0)
	xMax, yMax := min(xBmu+lim, s.size.Width-1), min(yBmu+lim, s.size.Height-1)

	for x := xMin; x <= xMax; x++ {
		for y := yMin; y <= yMax; y++ {
			r := s.neighborhood.Weight(x, y, xBmu, yBmu, radius)
			if r <= 0 {
				continue
			}
			for l, lay := range s.layers {
				node := lay.GetNode(x, y)
				for i := 0; i < lay.Columns(); i++ {
					d := data[l][i]
					if math.IsNaN(d) {
						continue
					}
					node[i] += alpha * r * (d - node[i])
				}
			}
		}
	}
}

func (s *Som) updateWeightsVI(bmuIdx int, data [][]float64, alpha, radius, lambda float64) {
	lim := s.neighborhood.MaxRadius(radius)
	if lim < 0 {
		lim = s.size.Width * s.size.Height
	}

	xBmu, yBmu := s.size.CoordsAt(bmuIdx)
	xMin, yMin := max(xBmu-lim, 0), max(yBmu-lim, 0)
	xMax, yMax := min(xBmu+lim, s.size.Width-1), min(yBmu+lim, s.size.Height-1)

	for x := xMin; x <= xMax; x++ {
		for y := yMin; y <= yMax; y++ {
			r := s.neighborhood.Weight(x, y, xBmu, yBmu, radius)
			if r <= 0 {
				continue
			}

			scale := 0.0
			if x != xBmu || y != yBmu {
				dataDist := s.nodeDistance(bmuIdx, s.size.IndexAt(x, y))
				mapDist := lambda * s.nodeMapDistance(xBmu, yBmu, x, y)
				scale = dataDist/mapDist - 1
				fmt.Fprintln(os.Stderr, "scale", scale, dataDist, mapDist, x, y)
				if dataDist > 100 {
					fmt.Fprintln(os.Stderr, "dataDist", dataDist, mapDist)
					fmt.Fprintln(os.Stderr, "node", s.layers[0].GetNode(x, y))
					panic("scale")
				}
			}

			for l, lay := range s.layers {
				bmu := lay.GetNodeAt(bmuIdx)
				node := lay.GetNode(x, y)
				for i := 0; i < lay.Columns(); i++ {
					d := data[l][i]
					if math.IsNaN(d) {
						continue
					}
					//node[i] += alpha * r * (d - node[i])
					//node[i] += alpha * r * ((d - bmu[i]) + (bmu[i] - node[i]))
					node[i] += alpha * r * ((d - bmu[i]) + (bmu[i]-node[i])*scale)

					if math.IsNaN(node[i]) || math.IsInf(node[i], 1) || math.IsInf(node[i], -1) {
						fmt.Fprintln(os.Stderr, "node value", node[i])
						fmt.Fprintln(os.Stderr, "bmu value", bmu[i], "node value", node[i], "data value", d, "scale", scale)
						fmt.Fprintln(os.Stderr, node)
						panic("NaN")
					}
				}
			}
		}
	}
}

func (s *Som) distance(data [][]float64, unit int) float64 {
	totalDist := 0.0
	for l, layer := range s.layers {
		node := layer.GetNodeAt(unit)
		dist := layer.Metric().Distance(node, data[l])
		totalDist += layer.Weight() * dist
	}
	return totalDist
}

func (s *Som) nodeDistance(unit1, unit2 int) float64 {
	totalDist := 0.0
	for _, layer := range s.layers {
		node1 := layer.GetNodeAt(unit1)
		node2 := layer.GetNodeAt(unit2)
		dist := layer.Metric().Distance(node1, node2)
		totalDist += layer.Weight() * dist
	}
	return totalDist
}

func (s *Som) nodeMapDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

func (s *Som) randomize(rng *rand.Rand) {
	for _, lay := range s.layers {
		data := lay.Data()
		for i := range data {
			data[i] = rng.Float64()
		}
	}
}

func (s *Som) Layers() []*layer.Layer {
	return s.layers
}

func (s *Som) UMatrix() [][]float64 {
	height := s.size.Height*2 - 1
	width := s.size.Width*2 - 1
	u := make([][]float64, height)

	for y := range u {
		u[y] = make([]float64, width)
		for x := range u[y] {
			u[y][x] = math.NaN()
		}
	}

	for x := 0; x < s.size.Width; x++ {
		for y := 0; y < s.size.Height; y++ {
			nodeHere := s.size.IndexAt(x, y)
			if x < s.size.Width-1 {
				nodeRight := s.size.IndexAt(x+1, y)
				u[y*2][x*2+1] = s.nodeDistance(nodeHere, nodeRight)
			}
			if y < s.size.Height-1 {
				nodeDown := s.size.IndexAt(x, y+1)
				u[y*2+1][x*2] = s.nodeDistance(nodeHere, nodeDown)
			}
		}
	}

	for x := 0; x < s.size.Width*2-1; x++ {
		for y := 0; y < s.size.Height*2-1; y++ {
			if x%2 != y%2 {
				continue
			}

			sum := 0.0
			cnt := 0

			if x > 0 {
				sum += u[y][x-1]
				cnt++
			}
			if x < width-1 {
				sum += u[y][x+1]
				cnt++
			}

			if y > 0 {
				sum += u[y-1][x]
				cnt++
			}
			if y < height-1 {
				sum += u[y+1][x]
				cnt++
			}

			u[y][x] = sum / float64(cnt)
		}
	}

	return u
}
