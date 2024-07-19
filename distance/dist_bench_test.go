package distance_test

import (
	"testing"

	"github.com/mlange-42/som/distance"
)

func benchmarkDistanceMetric(b *testing.B, d distance.Distance, dims int) {
	b.StopTimer()
	x, y := make([]float64, dims), make([]float64, dims)
	for i := 0; i < dims; i++ {
		f := float64(i) / float64(dims)
		x[i] = f
		y[i] = 1 - f
	}
	var dst float64
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		dst = d.Distance(x, y)
	}
	b.StopTimer()

	if dst != 0 {
		dst = 0
	}
	_ = dst
}

func BenchmarkSumOfSquaresDistance3(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.SumOfSquares{}, 3)
}

func BenchmarkSumOfSquaresDistance10(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.SumOfSquares{}, 10)
}

func BenchmarkEuclideanDistance3(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.Euclidean{}, 3)
}

func BenchmarkEuclideanDistance10(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.Euclidean{}, 10)
}

func BenchmarkManhattanDistance3(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.Manhattan{}, 3)
}

func BenchmarkManhattanDistance10(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.Manhattan{}, 10)
}

func BenchmarkHammingDistance3(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.Hamming{}, 3)
}

func BenchmarkHammingDistance10(b *testing.B) {
	benchmarkDistanceMetric(b, &distance.Hamming{}, 10)
}
