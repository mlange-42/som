package som

import (
	"container/heap"

	"gonum.org/v1/gonum/spatial/kdtree"
)

type nodeLocation struct {
	data      [][]float64
	nodeIndex int
}

func NewNodeLocation(index int, data [][]float64) *nodeLocation {
	return &nodeLocation{
		data:      data,
		nodeIndex: index,
	}
}

func (p *nodeLocation) GetIndex(d kdtree.Dim) (int, int) {
	dim := int(d)
	cum := 0
	for i, l := range p.data {
		c := cum
		cum += len(l)
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
	return p.data[lay][col] - q.data[lay][col]
}

// Dims returns the number of dimensions described by the receiver.
func (p nodeLocation) Dims() int {
	dims := 0
	for _, v := range p.data {
		dims += len(v)
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
	return loc.data[lay][col] < p.NodeLocations[j].data[lay][col]
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

// NDistKeeper keeps man number and distance
type NDistKeeper struct {
	kdtree.Heap
}

// NewNDistKeeper returns an NDistKeeper with the maximum value of the heap set to d.
func NewNDistKeeper(n int, d float64) *NDistKeeper {
	k := NDistKeeper{make(kdtree.Heap, 1, n)}
	k.Heap[0].Dist = d * d
	return &k
}

// Keep adds c to the heap if its distance is less than or equal to the max value of the heap.
func (k *NDistKeeper) Keep(c kdtree.ComparableDist) {
	if c.Dist <= k.Heap[0].Dist { // Favour later finds to displace sentinel.
		if len(k.Heap) == cap(k.Heap) {
			heap.Pop(k)
		}
		heap.Push(k, c)
	}
}
