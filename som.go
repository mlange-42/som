package som

import (
	"fmt"
	"math"
	"math/rand"
	"slices"

	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
	"github.com/mlange-42/som/table"
)

// SomConfig represents the configuration for a Self-Organizing Map (SOM).
// It defines the size of the map, the layers of data to be mapped, the neighborhood function,
// and the metric used to calculate distances on the map.
type SomConfig struct {
	Size         layer.Size                // Size of the SOM
	Layers       []*LayerDef               // Layer definitions
	Neighborhood neighborhood.Neighborhood // Neighborhood function of the SOM
	MapMetric    neighborhood.Metric       // Metric used to calculate distances on the map
	ViSomMetric  neighborhood.Metric       // Metric used to calculate distances on the map for ViSOM update
}

// PrepareTables reads the CSV data and creates a table for each layer defined in the SomConfig.
// If a categorical layer has no columns specified, it will attempt to read the class names for that layer
// and create a table from the classes. The created tables are returned in the same order as
// the layers in the SomConfig.
func (c *SomConfig) PrepareTables(reader table.Reader, ignoreLayers []string, updateNormalizers bool, keepOriginal bool) (normalized, raw []*table.Table, err error) {
	normalized = make([]*table.Table, len(c.Layers))
	raw = make([]*table.Table, len(c.Layers))

	ignoreFound := make([]bool, len(ignoreLayers))
	for i := range c.Layers {
		layer := c.Layers[i]

		if idx := slices.Index(ignoreLayers, layer.Name); idx >= 0 {
			ignoreFound[idx] = true
			continue
		}

		if layer.Categorical {
			tab, err := createCategoricalTable(reader, layer)
			if err != nil {
				return nil, nil, err
			}

			normalized[i] = tab

			err = keepTable(raw, i, tab, keepOriginal)
			if err != nil {
				return nil, nil, err
			}

			continue
		}

		if len(layer.Columns) == 0 {
			return nil, nil, fmt.Errorf("layer %s has no columns", layer.Name)
		}

		tab, err := reader.ReadColumns(layer.Columns)
		if err != nil {
			return nil, nil, err
		}

		err = keepTable(raw, i, tab, keepOriginal)
		if err != nil {
			return nil, nil, err
		}

		normalizeTable(tab, layer, updateNormalizers)
		normalized[i] = tab
	}

	for i, f := range ignoreFound {
		if !f {
			return nil, nil, fmt.Errorf("layer %s from ignore list not found in layers", ignoreLayers[i])
		}
	}

	if keepOriginal {
		return normalized, raw, nil
	}

	return normalized, nil, nil
}

func normalizeTable(tab *table.Table, layer *LayerDef, update bool) {
	if len(layer.Norm) != 0 {
		for j := range layer.Columns {
			if update {
				layer.Norm[j].Initialize(tab, j)
			}
			tab.NormalizeColumn(j, layer.Norm[j])
		}
	}
}

func createCategoricalTable(reader table.Reader, layer *LayerDef) (*table.Table, error) {
	classes, err := reader.ReadLabels(layer.Name)
	if err != nil {
		return nil, err
	}

	tab, err := conv.ClassesToTable(classes, layer.Columns, reader.NoData())
	if err != nil {
		return nil, err
	}
	layer.Columns = tab.ColumnNames()
	return tab, nil
}

func keepTable(list []*table.Table, idx int, tab *table.Table, keep bool) error {
	if !keep {
		return nil
	}
	rawTable, err := table.NewWithData(tab.ColumnNames(), append([]float64{}, tab.Data()...))
	if err != nil {
		return err
	}
	list[idx] = rawTable
	return nil
}

// LayerDef represents the configuration for a single layer in a Self-Organizing Map (SOM).
// It defines the name, columns, normalization, metric, weight, and whether the layer is categorical.
// If the layer has weights, it can also be initialized with the provided data.
//
// A weight value of 0.0 is interpreted as standard weight of 1.0.
// To get a weight of 0.0, give the weight field a negative value.
type LayerDef struct {
	Name        string            // Name of the layer
	Columns     []string          // Columns to use from the data
	Norm        []norm.Normalizer // Normalization functions for the columns
	Metric      distance.Distance // Distance metric to use for this layer
	Weight      float64           // Weight value for this layer (for multi-layer SOMs)
	Categorical bool              // Whether the layer contains categorical data
	Weights     []float64         // Pre-computed layer weights (if provided)
}

// Som represents a Self-Organizing Map (SOM) model.
type Som struct {
	size         layer.Size
	layers       []*layer.Layer
	neighborhood neighborhood.Neighborhood
	metric       neighborhood.Metric
	viSomMetric  neighborhood.Metric
}

// New creates a new Self-Organizing Map (SOM) instance based on the provided SomConfig.
// It initializes the layers of the SOM with the specified configurations.
// If the layer has pre-computed weights, they can be provided to initialize the layer.
// The function returns the created SOM instance and an error if any issues occur during
// the initialization.
func New(params *SomConfig) (*Som, error) {
	lay := make([]*layer.Layer, len(params.Layers))
	for i, l := range params.Layers {
		if len(l.Columns) == 0 {
			return nil, fmt.Errorf("layer %s has no columns", l.Name)
		}
		norm, err := checkAndFixLayerNorm(l)
		if err != nil {
			return nil, err
		}

		weight := l.Weight
		if weight == 0 {
			weight = 1
		} else if weight < 0 {
			weight = 0
		}
		metric := l.Metric
		if metric == nil {
			if l.Categorical {
				metric = &distance.Hamming{}
			} else {
				metric = &distance.Euclidean{}
			}
		}

		if len(l.Weights) == 0 {
			lay[i], err = layer.New(l.Name, l.Columns, norm, params.Size, metric, weight, l.Categorical)
		} else {
			lay[i], err = layer.NewWithData(l.Name, l.Columns, norm, params.Size, metric, weight, l.Categorical, l.Weights)
		}
		if err != nil {
			return nil, err
		}
	}
	return &Som{
		size:         params.Size,
		layers:       lay,
		neighborhood: params.Neighborhood,
		metric:       params.MapMetric,
		viSomMetric:  params.ViSomMetric,
	}, nil
}

func checkAndFixLayerNorm(l *LayerDef) ([]norm.Normalizer, error) {
	n := l.Norm
	if l.Categorical {
		if len(n) != len(l.Columns) && len(n) > 1 {
			return nil, fmt.Errorf("number of normalizers (%d) must match number of columns (%d) for layer %s", len(l.Norm), len(l.Columns), l.Name)
		}
		for _, nn := range n {
			if _, ok := nn.(*norm.Identity); !ok {
				return nil, fmt.Errorf("categorical layer %s must use identity normalizer", l.Name)
			}
		}
		if len(n) != len(l.Columns) {
			n = make([]norm.Normalizer, len(l.Columns))
			for j := range l.Columns {
				n[j] = &norm.Identity{}
			}
		}
	} else {
		if len(n) != len(l.Columns) {
			return nil, fmt.Errorf("number of normalizers (%d) must match number of columns (%d) for layer %s", len(l.Norm), len(l.Columns), l.Name)
		}
	}
	return n, nil
}

// Size returns the size of the Self-Organizing Map (SOM) instance.
func (s *Som) Size() *layer.Size {
	return &s.size
}

// Neighborhood returns the neighborhood function used by the Self-Organizing Map (SOM) instance.
func (s *Som) Neighborhood() neighborhood.Neighborhood {
	return s.neighborhood
}

// MapMetric returns the metric used to calculate distances between nodes in the Self-Organizing Map (SOM).
func (s *Som) MapMetric() neighborhood.Metric {
	return s.metric
}

// ViSomMetric returns the metric used to calculate distances between nodes in the ViSOM (Visualization Induced Self-Organizing Map).
func (s *Som) ViSomMetric() neighborhood.Metric {
	return s.viSomMetric
}

// Learn updates the weights of the Self-Organizing Map (SOM) based on the given input data.
// It calculates the Best Matching Unit (BMU) for the input data, then updates the weights
// of the nodes in the SOM based on the neighborhood function and learning rate.
// The function returns the distance between the input data and the BMU.
func (s *Som) Learn(data [][]float64, alpha, radius, lambda float64) float64 {
	bmuIdx, dist := s.GetBMU(data)

	s.updateWeights(bmuIdx, data, alpha, radius, lambda)

	return dist
}

// GetBMU finds the Best Matching Unit (BMU) for the given input data.
// It calculates the total distance between the input data and each node in the SOM,
// and returns the index of the node with the minimum total distance, along with that minimum distance.
func (s *Som) GetBMU(data [][]float64) (int, float64) {
	units := s.size.Nodes()

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

func (s *Som) GetBMU2(data [][]float64) (int, float64, int, float64) {
	units := s.size.Nodes()

	minDist := math.MaxFloat64
	minDist2 := math.MaxFloat64
	minIndex := -1
	minIndex2 := -1
	for i := 0; i < units; i++ {
		totalDist := s.distance(data, i)

		if totalDist < minDist {
			minDist2 = minDist
			minIndex2 = minIndex
			minDist = totalDist
			minIndex = i
		} else if totalDist < minDist2 {
			minDist2 = totalDist
			minIndex2 = i
		}
	}

	return minIndex, minDist, minIndex2, minDist2
}

func (s *Som) updateWeights(bmuIdx int, data [][]float64, alpha, radius, lambda float64) {
	lim := s.neighborhood.MaxRadius(radius)
	if lim < 0 {
		lim = s.size.Nodes()
	}

	xBmu, yBmu := s.size.Coords(bmuIdx)
	xMin, yMin := max(xBmu-lim, 0), max(yBmu-lim, 0)
	xMax, yMax := min(xBmu+lim, s.size.Width-1), min(yBmu+lim, s.size.Height-1)

	// update BMU, to use its new position in neighborhood updates
	s.updateNode(xBmu, yBmu, data, alpha)

	for x := xMin; x <= xMax; x++ {
		for y := yMin; y <= yMax; y++ {
			if x == xBmu && y == yBmu {
				// Skip BMU, already updated above
				continue
			}

			dist := s.metric.Distance(xBmu, yBmu, x, y)
			r := s.neighborhood.Weight(dist, radius)
			if r <= 0 {
				// Outside neighborhood, don't update
				continue
			}
			rate := r * alpha

			if lambda <= 0 {
				// Basic SOM
				s.updateNode(x, y, data, rate)
				continue
			}
			// ViSOM
			s.updateNodeViSom(bmuIdx, x, y, data, rate, lambda)
		}
	}
}

func (s *Som) updateNode(x, y int, data [][]float64, rate float64) {
	for l, lay := range s.layers {
		node := lay.GetNode(x, y)
		for i := 0; i < lay.Columns(); i++ {
			d := data[l][i]
			if math.IsNaN(d) {
				continue
			}
			node[i] += rate * (d - node[i])
		}
	}
}

func (s *Som) updateNodeViSom(bmuIdx, x, y int, data [][]float64, rate float64, lambda float64) {
	xBmu, yBmu := s.size.Coords(bmuIdx)
	nodeIdx := s.size.Index(x, y)
	d := s.nodeDistance(bmuIdx, nodeIdx)                   // distance in data space
	D := lambda * s.viSomMetric.Distance(xBmu, yBmu, x, y) // scaled distance in map space

	// scale = (d - D) / D = d/D - 1 (original formulation Yin 2002)
	scale := 0.0
	if d > 0 {
		scale = 1 - D/d // corrected formulation
	}

	for l, lay := range s.layers {
		bmu := lay.GetNodeAt(bmuIdx)
		node := lay.GetNodeAt(nodeIdx)
		for i := 0; i < lay.Columns(); i++ {
			d := data[l][i]
			if math.IsNaN(d) {
				continue
			}
			delta := (d - bmu[i]) + (bmu[i]-node[i])*scale

			node[i] += rate * delta
		}
	}
}

func (s *Som) decayWeights(center [][]float64, rate float64) {
	nodes := s.Size().Nodes()
	fac := 1.0 - rate

	for i := 0; i < nodes; i++ {
		for j, lay := range s.layers {
			node := lay.GetNodeAt(i)
			data := center[j]
			for k := 0; k < lay.Columns(); k++ {
				delta := node[k] - data[k]
				node[k] = data[k] + fac*delta
			}
		}
	}
}

func (s *Som) distance(data [][]float64, unit int) float64 {
	totalDist := 0.0
	for l, layer := range s.layers {
		if layer.Weight() == 0 || data[l] == nil {
			continue
		}
		node := layer.GetNodeAt(unit)
		dist := layer.Metric().Distance(node, data[l])
		totalDist += layer.Weight() * dist
	}
	return totalDist
}

func (s *Som) nodeDistance(unit1, unit2 int) float64 {
	totalDist := 0.0
	for _, layer := range s.layers {
		if layer.Weight() == 0 {
			continue
		}
		node1 := layer.GetNodeAt(unit1)
		node2 := layer.GetNodeAt(unit2)
		dist := layer.Metric().Distance(node1, node2)
		totalDist += layer.Weight() * dist
	}
	return totalDist
}

func (s *Som) Randomize(rng *rand.Rand) {
	for _, lay := range s.layers {
		data := lay.Weights()
		for i := range data {
			data[i] = rng.Float64() * 0.25
		}
	}
}

// Layers returns the layers of the Self-Organizing Map.
func (s *Som) Layers() []*layer.Layer {
	return s.layers
}

// UMatrix computes the U-Matrix for the Self-Organizing Map. The U-Matrix
// visualizes the distances between neighboring nodes in the map, which can
// be used to identify cluster boundaries. The returned matrix has double
// the dimensions of the original map, with the values representing the
// distances between nodes and their neighbors.
//
// If fill is true, cells that don't correspond to a link, but to a node or an "empty space"
// are filled with the average of the surrounding links.
func (s *Som) UMatrix(fill bool) [][]float64 {
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
			nodeHere := s.size.Index(x, y)
			if x < s.size.Width-1 {
				nodeRight := s.size.Index(x+1, y)
				u[y*2][x*2+1] = s.nodeDistance(nodeHere, nodeRight)
			}
			if y < s.size.Height-1 {
				nodeDown := s.size.Index(x, y+1)
				u[y*2+1][x*2] = s.nodeDistance(nodeHere, nodeDown)
			}
		}
	}

	if !fill {
		return u
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
