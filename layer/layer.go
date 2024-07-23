package layer

import (
	"fmt"
	"slices"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/norm"
)

// Size represents the width and height of a 2D layer or grid.
type Size struct {
	Width  int
	Height int
}

// Coords returns the (x, y) coordinates of the node at the given index.
func (s *Size) Coords(idx int) (int, int) {
	return idx / s.Height, idx % s.Height
}

func (s *Size) Index(x, y int) int {
	return y + x*s.Height
}

func (s *Size) Nodes() int {
	return s.Width * s.Height
}

// Layer represents a layer of data in a Self-organizing Map.
type Layer struct {
	name        string            // The name of the layer
	columns     []string          // The names of the columns in the layer
	norm        []norm.Normalizer // The normalizers for the layer
	size        Size              // The width and height of the layer
	weight      float64           // The weight of the layer
	metric      distance.Distance // The distance metric for the layer
	weights     []float64         // The weight values for the layer
	categorical bool              // Whether the layer is categorical or continuous
}

// New creates a new Layer.
func New(name string, columns []string, normalizers []norm.Normalizer, size Size, metric distance.Distance, weight float64, categorical bool) (*Layer, error) {
	return NewWithData(
		name, columns, normalizers,
		size, metric, weight, categorical,
		make([]float64, size.Width*size.Height*len(columns)),
	)
}

// NewWithData creates a new Layer with the given initial data.
func NewWithData(name string, columns []string, normalizers []norm.Normalizer, size Size, metric distance.Distance, weight float64, categorical bool, data []float64) (*Layer, error) {
	if len(data) != size.Width*size.Height*len(columns) {
		return nil, fmt.Errorf("data length (%d) does not match layer size (%d)", len(data), size.Width*size.Height*len(columns))
	}
	if len(normalizers) == 0 {
		normalizers = make([]norm.Normalizer, len(columns))
		for i := range columns {
			normalizers[i] = &norm.Identity{}
		}
	}
	if len(normalizers) != len(columns) {
		panic(fmt.Sprintf("invalid number of normalizers: expected %d, got %d", len(columns), len(normalizers)))
	}
	return &Layer{
		name:        name,
		columns:     columns,
		norm:        normalizers,
		metric:      metric,
		weight:      weight,
		size:        size,
		weights:     data,
		categorical: categorical,
	}, nil
}

// Name returns the name of the Layer.
func (l *Layer) Name() string {
	return l.name
}

// Weights returns the weight values for the layer.
func (l *Layer) Weights() []float64 {
	return l.weights
}

// Metric returns the distance metric used by the Layer.
func (l *Layer) Metric() distance.Distance {
	return l.metric
}

// Weight returns the weight value of the Layer.
// The weight value is used in distance calculations for multi-layer SOMs.
func (l *Layer) Weight() float64 {
	return l.weight
}

func (l *Layer) nodeIndex(x, y int) int {
	return (y + x*l.size.Height) * len(l.columns)
}

func (l *Layer) nodeIndexAt(idx int) int {
	return idx * len(l.columns)
}

func (l *Layer) index(x, y, col int) int {
	return (y+x*l.size.Height)*len(l.columns) + col
}

func (l *Layer) indexAt(idx, col int) int {
	return idx*len(l.columns) + col
}

// Column returns the index of the column with the given name in the Layer.
// If the column is not found, it returns -1.
func (l *Layer) Column(col string) int {
	return slices.Index(l.columns, col)
}

// ColumnNames returns the names of the columns in the Layer.
func (l *Layer) ColumnNames() []string {
	return l.columns
}

// Columns returns the number of columns in the Layer.
func (l *Layer) Columns() int {
	return len(l.columns)
}

// Nodes returns the total number of nodes in the Layer.
func (l *Layer) Nodes() int {
	return l.size.Nodes()
}

// Normalizers returns the list of normalizers used by the Layer.
func (l *Layer) Normalizers() []norm.Normalizer {
	return l.norm
}

// Get returns the value at the specified column and coordinate in the Layer.
func (l *Layer) Get(x, y, col int) float64 {
	return l.weights[l.index(x, y, col)]
}

// GetAt returns the value at the specified column and node index in the Layer.
func (l *Layer) GetAt(idx, col int) float64 {
	return l.weights[l.indexAt(idx, col)]
}

// Set sets the value at the specified column and coordinate in the Layer.
func (l *Layer) Set(x, y, col int, value float64) {
	l.weights[l.index(x, y, col)] = value
}

// SetAt sets the value at the specified column and node index in the Layer.
func (l *Layer) SetAt(idx, col int, value float64) {
	l.weights[l.indexAt(idx, col)] = value
}

// GetNode returns a slice of float64 values representing the data for the node
// at the specified (x, y) coordinates in the Layer. The slice contains the
// values for each column in the Layer, in the same order as the columns slice.
func (l *Layer) GetNode(x, y int) []float64 {
	idx := l.nodeIndex(x, y)
	return l.weights[idx : idx+len(l.columns)]
}

// GetNodeAt returns a slice of float64 values representing the data for the node
// at the index in the Layer. The slice contains the
// values for each column in the Layer, in the same order as the columns slice.
func (l *Layer) GetNodeAt(idx int) []float64 {
	idx2 := l.nodeIndexAt(idx)
	return l.weights[idx2 : idx2+len(l.columns)]
}

// CoordsAt returns the (x, y) coordinates of the node at the specified index.
func (l *Layer) CoordsAt(idx int) (int, int) {
	return l.size.Coords(idx)
}

// IsCategorical returns whether the Layer contains categorical data.
func (l *Layer) IsCategorical() bool {
	return l.categorical
}

// ColumnMatrix returns a 2D slice of float64 values representing the data for the
// specified column in the Layer. Each inner slice represents a row in the Layer,
// and the values in the slice correspond to the values in that column.
func (l *Layer) ColumnMatrix(col int) [][]float64 {
	data := make([][]float64, l.size.Height)
	for y := 0; y < l.size.Height; y++ {
		row := make([]float64, l.size.Width)
		for x := 0; x < l.size.Width; x++ {
			row[x] = l.Get(x, y, col)
		}
		data[y] = row
	}
	return data
}

// DeNormalize all weight vectors in place.
// Applies the inverse of each column's normalizer to the column's values.
func (l *Layer) DeNormalize() {
	nodes := l.Nodes()
	cols := len(l.columns)
	for i := 0; i < nodes; i++ {
		for col := 0; col < cols; col++ {
			l.SetAt(i, col, l.norm[col].DeNormalize(l.GetAt(i, col)))
		}
	}
}
