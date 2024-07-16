package som

import "slices"

// Size represents the width and height of a 2D layer or grid.
type Size struct {
	Width  int
	Height int
}

// Layer represents a layer of data in a Self-organizing Map.
type Layer struct {
	columns []string  // The names of the columns in the layer
	size    Size      // The width and height of the layer
	data    []float64 // The data values for the layer
}

// NewLayer creates a new Layer with the given columns and size.
func NewLayer(columns []string, size Size) Layer {
	return Layer{
		columns: columns,
		size:    size,
		data:    make([]float64, size.Width*size.Height*len(columns)),
	}
}

func (l *Layer) nodeIndex(x, y int) int {
	return (y + x*l.size.Height) * len(l.columns)
}

func (l *Layer) index(x, y, col int) int {
	return (y+x*l.size.Height)*len(l.columns) + col
}

// Column returns the index of the column with the given name in the Layer.
// If the column is not found, it returns -1.
func (l *Layer) Column(col string) int {
	return slices.Index(l.columns, col)
}

// Get returns the value at the specified column and coordinate in the Layer.
func (l *Layer) Get(x, y, col int) float64 {
	return l.data[l.index(x, y, col)]
}

// GetNode returns a slice of float64 values representing the data for the node
// at the specified (x, y) coordinates in the Layer. The slice contains the
// values for each column in the Layer, in the same order as the columns slice.
func (l *Layer) GetNode(x, y int) []float64 {
	idx := l.nodeIndex(x, y)
	return l.data[idx : idx+len(l.columns)]
}
