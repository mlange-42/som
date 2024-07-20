package neighborhood

import "math"

var metrics = map[string]Metric{}

func init() {
	m := []Metric{
		&Euclidean{},
		&Manhattan{},
		&Chebyshev{},
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

type Metric interface {
	Name() string
	Distance(x1, y1, x2, y2 int) float64
}

type Euclidean struct{}

func (e *Euclidean) Name() string {
	return "euclidean"
}

func (e *Euclidean) Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

type Manhattan struct{}

func (m *Manhattan) Name() string {
	return "manhattan"
}

func (m *Manhattan) Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Abs(dx) + math.Abs(dy)
}

type Chebyshev struct{}

func (c *Chebyshev) Name() string {
	return "chebyshev"
}

func (c *Chebyshev) Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Max(math.Abs(dx), math.Abs(dy))
}
