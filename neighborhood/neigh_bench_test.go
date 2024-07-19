package neighborhood_test

import (
	"testing"

	"github.com/mlange-42/som/neighborhood"
)

func benchmarkNeighborhood(b *testing.B, g neighborhood.Neighborhood) {
	b.StopTimer()

	x1, y1, x2, y2 := 1, 0, 0, -1
	radius := 3.0

	var w float64

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w = g.Weight(x1, y1, x2, y2, radius)
	}
	b.StopTimer()

	if w != 0 {
		w = 0
	}
	_ = w
}

func BenchmarkGaussianWeight(b *testing.B) {
	benchmarkNeighborhood(b, &neighborhood.Gaussian{})
}

func BenchmarkCutGaussianWeight(b *testing.B) {
	benchmarkNeighborhood(b, &neighborhood.CutGaussian{})
}

func BenchmarkLinearWeight(b *testing.B) {
	benchmarkNeighborhood(b, &neighborhood.Linear{})
}

func BenchmarkBoxWeight(b *testing.B) {
	benchmarkNeighborhood(b, &neighborhood.Box{})
}
