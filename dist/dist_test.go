package dist

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumOfSquaresDistance(t *testing.T) {
	d := &SumOfSquares{}
	tests := []struct {
		name     string
		x        []float64
		y        []float64
		expected float64
	}{
		{"Zero vectors", []float64{0, 0, 0}, []float64{0, 0, 0}, 0},
		{"All same", []float64{3, 2, 1}, []float64{3, 2, 1}, 0},
		{"Positive vectors", []float64{1, 2, 3}, []float64{4, 5, 6}, 27},
		{"Mixed vectors", []float64{-1, 2, -3}, []float64{1, -2, 3}, 56},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.Distance(tt.x, tt.y)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEuclideanDistance(t *testing.T) {
	d := &Euclidean{}
	tests := []struct {
		name     string
		x        []float64
		y        []float64
		expected float64
	}{
		{"Zero vectors", []float64{0, 0, 0}, []float64{0, 0, 0}, 0},
		{"All same", []float64{3, 2, 1}, []float64{3, 2, 1}, 0},
		{"Positive vectors", []float64{1, 1, 1}, []float64{4, 5, 6}, math.Sqrt(3*3 + 4*4 + 5*5)},
		{"Mixed vectors", []float64{-1, 2, -3}, []float64{1, -2, 3}, math.Sqrt(2*2 + 4*4 + 6*6)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.Distance(tt.x, tt.y)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManhattanDistance(t *testing.T) {
	d := &Manhattan{}
	tests := []struct {
		name     string
		x        []float64
		y        []float64
		expected float64
	}{
		{"Zero vectors", []float64{0, 0, 0}, []float64{0, 0, 0}, 0},
		{"All same", []float64{3, 2, 1}, []float64{3, 2, 1}, 0},
		{"Positive vectors", []float64{1, 1, 1}, []float64{4, 5, 6}, 3 + 4 + 5},
		{"Mixed vectors", []float64{-1, 2, -3}, []float64{1, -2, 3}, 2 + 4 + 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.Distance(tt.x, tt.y)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHammingDistance(t *testing.T) {
	d := &Hamming{}
	tests := []struct {
		name     string
		x        []float64
		y        []float64
		expected float64
	}{
		{"Zero vectors", []float64{0, 0, 0}, []float64{0, 0, 0}, 0},
		{"All same", []float64{1, 0, 0, 1, 1}, []float64{1, 0, 0, 1, 1}, 0},
		{"All different", []float64{1, 0, 0, 1, 1}, []float64{0, 1, 1, 0, 0}, 1},
		{"Some different", []float64{1, 0, 0, 1, 1}, []float64{1, 1, 0, 1, 0}, 0.4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.Distance(tt.x, tt.y)
			assert.Equal(t, tt.expected, result)
		})
	}
}
