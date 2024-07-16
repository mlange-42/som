package som

import "slices"

type Size struct {
	Width  int
	Height int
}

type Layer struct {
	columns []string
	size    Size
	data    []float64
}

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

func (l *Layer) Column(col string) int {
	return slices.Index(l.columns, col)
}

func (l *Layer) Get(x, y, col int) float64 {
	return l.data[l.index(x, y, col)]
}

func (l *Layer) GetNode(x, y int) []float64 {
	idx := l.nodeIndex(x, y)
	return l.data[idx : idx+len(l.columns)]
}
