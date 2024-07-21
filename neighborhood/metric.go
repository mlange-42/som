package neighborhood

import "math"

var metrics = map[string]Metric{}

func init() {
	m := []Metric{
		&EuclideanMetric{},
		&ManhattanMetric{},
		&ChebyshevMetric{},
	}
	for _, v := range m {
		if _, ok := metrics[v.Name()]; ok {
			panic("duplicate metric name: " + v.Name())
		}
		metrics[v.Name()] = v
	}
}

func GetMetric(name string) (Metric, bool) {
	m, ok := metrics[name]
	return m, ok
}

// Metric is an interface that defines a distance metric in map space, i.e. between SOM nodes.
// The Name method returns the name of the metric.
// The Distance method calculates the distance between two points (x1, y1) and (x2, y2).
type Metric interface {
	Name() string
	Distance(x1, y1, x2, y2 int) float64
}

// EuclideanMetric implements [Metric] for the Euclidean distance.
type EuclideanMetric struct{}

func (e *EuclideanMetric) Name() string {
	return "euclidean"
}

func (e *EuclideanMetric) Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

// ManhattanMetric implements [Metric] for the Manhattan distance.
type ManhattanMetric struct{}

func (m *ManhattanMetric) Name() string {
	return "manhattan"
}

func (m *ManhattanMetric) Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Abs(dx) + math.Abs(dy)
}

// ChebyshevMetric implements [Metric] for the Chebyshev distance.
type ChebyshevMetric struct{}

func (c *ChebyshevMetric) Name() string {
	return "chebyshev"
}

func (c *ChebyshevMetric) Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Max(math.Abs(dx), math.Abs(dy))
}
