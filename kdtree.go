package som

import (
	"gonum.org/v1/gonum/spatial/kdtree"
)

type nodeLocation struct {
	som       *Som
	nodeIndex int
}

func NewNodeLocation(som *Som, index int) *nodeLocation {
	return &nodeLocation{
		som:       som,
		nodeIndex: index,
	}
}

func (p *nodeLocation) GetIndex(d kdtree.Dim) (int, int) {
	dim := int(d)
	cum := 0
	for i, l := range p.som.layers {
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
func (p nodeLocation) Compare(c kdtree.Comparable, d kdtree.Dim) float64 {
	q := c.(nodeLocation)
	lay, col := p.GetIndex(d)
	l := p.som.layers[lay]
	return l.GetAt(p.nodeIndex, col) - l.GetAt(q.nodeIndex, col)
}

// Dims returns the number of dimensions described by the receiver.
func (p nodeLocation) Dims() int {
	dims := 0
	for _, l := range p.som.layers {
		dims += l.Columns()
	}
	return dims
}

// Distance returns the squared Euclidean distance between c and the receiver. The
// concrete type of c must be EntityLocation.
func (p nodeLocation) Distance(c kdtree.Comparable) float64 {
	q := c.(nodeLocation)
	_ = q
	return 0
}

type NodeLocations []nodeLocation

func (p NodeLocations) Index(i int) kdtree.Comparable { return p[i] }
func (p NodeLocations) Len() int                      { return len(p) }
func (p NodeLocations) Pivot(d kdtree.Dim) int {
	return plane{NodeLocations: p, Dim: d}.Pivot()
}
func (p NodeLocations) Slice(start, end int) kdtree.Interface { return p[start:end] }

// plane is a wrapping type that allows a Points type be pivoted on a dimension.
// The Pivot method of Plane uses MedianOfRandoms sampling at most 100 elements
// to find a pivot element.
type plane struct {
	kdtree.Dim
	NodeLocations
}

// randoms is the maximum number of random values to sample for calculation of
// median of random elements.
const randoms = 100

// Less comparison
func (p plane) Less(i, j int) bool {
	loc := p.NodeLocations[i]
	lay, col := loc.GetIndex(p.Dim)
	l := loc.som.layers[lay]
	return l.GetAt(loc.nodeIndex, col) < l.GetAt(p.NodeLocations[j].nodeIndex, col)
}

// Pivot TreePlane
func (p plane) Pivot() int { return kdtree.Partition(p, kdtree.MedianOfRandoms(p, randoms)) }

// Slice TreePlane
func (p plane) Slice(start, end int) kdtree.SortSlicer {
	p.NodeLocations = p.NodeLocations[start:end]
	return p
}

// Swap TreePlane
func (p plane) Swap(i, j int) {
	p.NodeLocations[i], p.NodeLocations[j] = p.NodeLocations[j], p.NodeLocations[i]
}
