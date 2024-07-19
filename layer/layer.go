package layer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/norm"
)

// Size represents the width and height of a 2D layer or grid.
type Size struct {
	Width  int
	Height int
}

// CoordsAt returns the (x, y) coordinates of the node at the given index.
func (s *Size) CoordsAt(idx int) (int, int) {
	return idx / s.Height, idx % s.Height
}

// Layer represents a layer of data in a Self-organizing Map.
type Layer struct {
	name        string            // The name of the layer
	columns     []string          // The names of the columns in the layer
	norm        []norm.Normalizer // The normalizers for the layer
	size        Size              // The width and height of the layer
	weight      float64           // The weight of the layer
	metric      distance.Distance // The distance metric for the layer
	data        []float64         // The data values for the layer
	categorical bool              // Whether the layer is categorical or continuous
}

// New creates a new Layer with the given columns and size.
func New(name string, columns []string, normalizers []norm.Normalizer, size Size, metric distance.Distance, weight float64, categorical bool) (*Layer, error) {
	return NewWithData(
		name, columns, normalizers,
		size, metric, weight, categorical,
		make([]float64, size.Width*size.Height*len(columns)),
	)
}

func NewWithData(name string, columns []string, normalizers []norm.Normalizer, size Size, metric distance.Distance, weight float64, categorical bool, data []float64) (*Layer, error) {
	if len(data) != size.Width*size.Height*len(columns) {
		return nil, fmt.Errorf("data length (%d) does not match layer size (%d)", len(data), size.Width*size.Height*len(columns))
	}
	if len(normalizers) == 0 {
		normalizers = make([]norm.Normalizer, len(columns))
		for i := range columns {
			normalizers[i] = &norm.None{}
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
		data:        data,
		categorical: categorical,
	}, nil
}

func (l *Layer) Name() string {
	return l.name
}

func (l *Layer) Data() []float64 {
	return l.data
}

func (l *Layer) Metric() distance.Distance {
	return l.metric
}

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

func (l *Layer) ColumnNames() []string {
	return l.columns
}

func (l *Layer) Columns() int {
	return len(l.columns)
}

func (l *Layer) Nodes() int {
	return l.size.Width * l.size.Height
}

func (l *Layer) Normalizers() []norm.Normalizer {
	return l.norm
}

// Get returns the value at the specified column and coordinate in the Layer.
func (l *Layer) Get(x, y, col int) float64 {
	return l.data[l.index(x, y, col)]
}

// GetAt returns the value at the specified column and node index in the Layer.
func (l *Layer) GetAt(idx, col int) float64 {
	return l.data[l.indexAt(idx, col)]
}

func (l *Layer) Set(x, y, col int, value float64) {
	l.data[l.index(x, y, col)] = value
}

func (l *Layer) SetAt(idx, col int, value float64) {
	l.data[l.indexAt(idx, col)] = value
}

// GetNode returns a slice of float64 values representing the data for the node
// at the specified (x, y) coordinates in the Layer. The slice contains the
// values for each column in the Layer, in the same order as the columns slice.
func (l *Layer) GetNode(x, y int) []float64 {
	idx := l.nodeIndex(x, y)
	return l.data[idx : idx+len(l.columns)]
}

// GetNodeAt returns a slice of float64 values representing the data for the node
// at the index in the Layer. The slice contains the
// values for each column in the Layer, in the same order as the columns slice.
func (l *Layer) GetNodeAt(idx int) []float64 {
	idx2 := l.nodeIndexAt(idx)
	return l.data[idx2 : idx2+len(l.columns)]
}

func (l *Layer) CoordsAt(idx int) (int, int) {
	return l.size.CoordsAt(idx)
}

func (l *Layer) IsCategorical() bool {
	return l.categorical
}

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

func (l *Layer) ToCSV(sep rune) string {
	b := strings.Builder{}
	s := string(sep)
	cols := l.ColumnNames()

	b.WriteString(fmt.Sprintf("id%scol%srow%s", s, s, s))
	for i, col := range cols {
		b.WriteString(col)
		if i < len(cols)-1 {
			b.WriteRune(sep)
		}
	}
	b.WriteRune('\n')
	nodes := l.Nodes()
	for i := 0; i < nodes; i++ {
		col, row := l.CoordsAt(i)
		b.WriteString(fmt.Sprintf("%d%s%d%s%d%s", i, s, col, s, row, s))
		for j := 0; j < len(cols); j++ {
			b.WriteString(fmt.Sprintf("%f", l.GetAt(i, j)))
			if j < len(cols)-1 {
				b.WriteRune(sep)
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (l *Layer) DeNormalize() {
	nodes := l.Nodes()
	cols := len(l.columns)
	for i := 0; i < nodes; i++ {
		for col := 0; col < cols; col++ {
			l.SetAt(i, col, l.norm[col].DeNormalize(l.GetAt(i, col)))
		}
	}
}
