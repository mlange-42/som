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
	Distance(x, y []float64) float64
}

type SumOfSquares struct{}

func (d *SumOfSquares) Name() string {
	return "sumofsquares"
}

func (d *SumOfSquares) Distance(x, y []float64) float64 {
	var sum float64
	for i := range x {
		if math.IsNaN(y[i]) {
			continue
		}
		d := x[i] - y[i]
		sum += d * d
	}
	return sum
}

type Euclidean struct{}

func (d *Euclidean) Name() string {
	return "euclidean"
}

func (d *Euclidean) Distance(x, y []float64) float64 {
	var sum float64
	for i := range x {
		if math.IsNaN(y[i]) {
			continue
		}
		d := x[i] - y[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

type Manhattan struct{}

func (d *Manhattan) Name() string {
	return "manhattan"
}

func (d *Manhattan) Distance(x, y []float64) float64 {
	var sum float64
	for i := range x {
		if math.IsNaN(y[i]) {
			continue
		}
		sum += math.Abs(x[i] - y[i])
	}
	return sum
}

type Hamming struct{}

func (d *Hamming) Name() string {
	return "hamming"
}

func (d *Hamming) Distance(x, y []float64) float64 {
	var sum float64
	for i := range x {
		if math.IsNaN(y[i]) {
			continue
		}
		if (x[i] < 0.5) != (y[i] < 0.5) {
			sum++
		}
	}
	return sum / float64(len(x))
}
