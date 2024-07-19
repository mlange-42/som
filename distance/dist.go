package distance

import (
	"math"
)

var metrics = map[string]Distance{}

func init() {
	m := []Distance{
		&SumOfSquares{},
		&Euclidean{},
		&Manhattan{},
		&Hamming{},
	}
	for _, v := range m {
		if _, ok := metrics[v.Name()]; ok {
			panic("duplicate metric name: " + v.Name())
		}
		metrics[v.Name()] = v
	}
}

func GetMetric(name string) (Distance, bool) {
	d, ok := metrics[name]
	return d, ok
}

type Distance interface {
	Name() string
	Distance(node, data []float64) float64
}

type SumOfSquares struct{}

func (d *SumOfSquares) Name() string {
	return "sumofsquares"
}

func (d *SumOfSquares) Distance(node, data []float64) float64 {
	var sum float64
	for i := range node {
		if math.IsNaN(data[i]) {
			continue
		}
		d := node[i] - data[i]
		sum += d * d
	}
	return sum
}

type Euclidean struct{}

func (d *Euclidean) Name() string {
	return "euclidean"
}

func (d *Euclidean) Distance(node, data []float64) float64 {
	var sum float64
	for i := range node {
		if math.IsNaN(data[i]) {
			continue
		}
		d := node[i] - data[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

type Manhattan struct{}

func (d *Manhattan) Name() string {
	return "manhattan"
}

func (d *Manhattan) Distance(node, data []float64) float64 {
	var sum float64
	for i := range node {
		if math.IsNaN(data[i]) {
			continue
		}
		sum += math.Abs(node[i] - data[i])
	}
	return sum
}

type Hamming struct{}

func (d *Hamming) Name() string {
	return "hamming"
}

func (d *Hamming) Distance(node, data []float64) float64 {
	var sum float64
	for i := range node {
		if math.IsNaN(data[i]) {
			continue
		}
		if (node[i] < 0.5) != (data[i] < 0.5) {
			sum++
		}
	}
	return sum / float64(len(node))
}
