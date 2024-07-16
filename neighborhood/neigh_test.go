package neighborhood_test

import (
	"math"
	"testing"

	"github.com/mlange-42/som/neighborhood"
	"github.com/stretchr/testify/assert"
)

func TestGaussian(t *testing.T) {
	g := &neighborhood.Gaussian{Sigma: 1.0}

	tests := []struct {
		name     string
		x1, y1   int
		x2, y2   int
		decay    float64
		expected float64
	}{
		{"Same point", 0, 0, 0, 0, 1.0, 1.0},
		{"Unit distance", 1, 0, 0, 0, 1.0, math.Exp(-0.5)},
		{"Diagonal distance", 1, 1, 0, 0, 1.0, math.Exp(-1)},
		{"With decay", 1, 1, 0, 0, 0.5, math.Exp(-1 / 0.25)},
		{"Large distance", 10, 10, 0, 0, 1.0, math.Exp(-100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := g.Weight(tt.x1, tt.y1, tt.x2, tt.y2, tt.decay)
			assert.InDelta(t, tt.expected, result, 1e-9)
		})
	}
}

func TestCutGaussian(t *testing.T) {
	g := &neighborhood.CutGaussian{Sigma: 1.0}

	tests := []struct {
		name     string
		x1, y1   int
		x2, y2   int
		decay    float64
		expected float64
	}{
		{"Same point", 0, 0, 0, 0, 1.0, 1.0},
		{"Unit distance", 1, 0, 0, 0, 1.0, math.Exp(-0.5)},
		{"Diagonal distance", 1, 1, 0, 0, 1.0, math.Exp(-1)},
		{"With decay", 1, 1, 0, 0, 0.5, math.Exp(-1 / 0.25)},
		{"Large distance", 10, 10, 0, 0, 1.0, math.Exp(-100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := g.Weight(tt.x1, tt.y1, tt.x2, tt.y2, tt.decay)
			assert.InDelta(t, tt.expected, result, 1e-9)
		})
	}
}
func TestBox(t *testing.T) {
	b := &neighborhood.Box{Size: 2.0}

	tests := []struct {
		name     string
		x1, y1   int
		x2, y2   int
		decay    float64
		expected float64
	}{
		{"Inside box", 1, 1, 0, 0, 1.0, 1.0},
		{"On box edge", 2, 0, 0, 0, 1.0, 1.0},
		{"Outside box", 3, 0, 0, 0, 1.0, 0.0},
		{"Diagonal inside", 1, 1, 0, 0, 1.0, 1.0},
		{"Diagonal outside", 2, 2, 0, 0, 1.0, 0.0},
		{"With decay inside", 1, 0, 0, 0, 0.5, 1.0},
		{"With decay outside", 2, 0, 0, 0, 0.5, 0.0},
		{"Zero distance", 0, 0, 0, 0, 1.0, 1.0},
		{"Negative coordinates", -1, -1, 0, 0, 1.0, 1.0},
		{"Large distance", 100, 100, 0, 0, 1.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.Weight(tt.x1, tt.y1, tt.x2, tt.y2, tt.decay)
			assert.Equal(t, tt.expected, result)
		})
	}
}
