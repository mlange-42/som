package som

import (
	"gonum.org/v1/gonum/spatial/kdtree"
)

type nodeLocation struct {
	Som       *Som
	NodeIndex int
	Data      [][]float64
}

func newNodeLocation(som *Som, nodeIndex int) nodeLocation {
	return nodeLocation{
		Som:       som,
		NodeIndex: nodeIndex,
	}
}

func newDataLocation(som *Som, data [][]float64) nodeLocation {
	return nodeLocation{
		Som:       som,
		NodeIndex: -1,
		Data:      data,
	}
}

func (p *nodeLocation) Get(lay, col int) float64 {
	if p.NodeIndex < 0 {
		return p.Data[lay][col]
	}
	l := p.Som.layers[lay]
	return l.GetAt(p.NodeIndex, col)
}

func (p *nodeLocation) GetIndex(d kdtree.Dim) (int, int) {
	dim := int(d)
	cum := 0
	for i, l := range p.Som.layers {
		c := cum
		cum += l.Columns()
		if cum > dim {
			idx := dim - c
			return i, idx
		}
	}
	panic("invalid dimension index")
}

// Compare returns the signed distance of p from the plane passing through c and
// perpendicular to the dimension d. The concrete type of c must be EntityLocation.
func (p *nodeLocation) Compare(c kdtree.Comparable, d kdtree.Dim) float64 {
	q := c.(*nodeLocation)
	lay, col := p.GetIndex(d)
	return p.Get(lay, col) - q.Get(lay, col)
}

// Dims returns the number of dimensions described by the receiver.
func (p *nodeLocation) Dims() int {
	dims := 0
	for _, l := range p.Som.layers {
		dims += l.Columns()
	}
	return dims
}

// Distance returns the squared Euclidean distance between c and the receiver. The
// concrete type of c must be EntityLocation.
func (p *nodeLocation) Distance(c kdtree.Comparable) float64 {
	q := c.(*nodeLocation)
	if p.NodeIndex < 0 && q.NodeIndex < 0 {
		panic("cannot compute distance between data points, at least one SOM node is required")
	}
	if p.NodeIndex >= 0 && q.NodeIndex >= 0 {
		return p.Som.nodeDistance(p.NodeIndex, q.NodeIndex)
	}
	if q.NodeIndex < 0 {
		p, q = q, p
	}
	return p.Som.dataDistance(p.Data, q.NodeIndex)
}

type nodeLocations []nodeLocation

func newNodeLocations(som *Som) nodeLocations {
	locs := make([]nodeLocation, som.Size().Nodes())
	for i := range locs {
		locs[i] = newNodeLocation(som, i)
	}
	return locs
}

func (p nodeLocations) Index(i int) kdtree.Comparable { return &p[i] }
func (p nodeLocations) Len() int                      { return len(p) }
func (p nodeLocations) Pivot(d kdtree.Dim) int {
	return plane{nodeLocations: p, Dim: d}.Pivot()
}
func (p nodeLocations) Slice(start, end int) kdtree.Interface { return p[start:end] }

// plane is a wrapping type that allows a Points type be pivoted on a dimension.
// The Pivot method of Plane uses MedianOfRandoms sampling at most 100 elements
// to find a pivot element.
type plane struct {
	kdtree.Dim
	nodeLocations
}

// randoms is the maximum number of random values to sample for calculation of
// median of random elements.
const randoms = 100

// Less comparison
func (p plane) Less(i, j int) bool {
	loc := p.nodeLocations[i]
	lay, col := loc.GetIndex(p.Dim)
	return loc.Get(lay, col) < p.nodeLocations[j].Get(lay, col)
}

// Pivot TreePlane
func (p plane) Pivot() int { return kdtree.Partition(p, kdtree.MedianOfRandoms(p, randoms)) }

// Slice TreePlane
func (p plane) Slice(start, end int) kdtree.SortSlicer {
	p.nodeLocations = p.nodeLocations[start:end]
	return p
}

// Swap TreePlane
func (p plane) Swap(i, j int) {
	p.nodeLocations[i], p.nodeLocations[j] = p.nodeLocations[j], p.nodeLocations[i]
}
