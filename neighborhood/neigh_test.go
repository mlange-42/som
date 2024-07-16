package neighborhood_test

import (
	"math"
	"testing"

	"github.com/mlange-42/som/neighborhood"
)

func TestGaussianWeight(t *testing.T) {
	g := &neighborhood.Gaussian{}
	g2 := &neighborhood.CutGaussian{}

	tests := []struct {
		name   string
		x1, y1 int
		x2, y2 int
		radius float64
		want   float64
	}{
		{
			name: "Same point",
			x1:   0, y1: 0,
			x2: 0, y2: 0,
			radius: 1.0,
			want:   1.0,
		},
		{
			name: "Points at radius distance",
			x1:   0, y1: 0,
			x2: 1, y2: 0,
			radius: 1.0,
			want:   math.Exp(-0.5),
		},
		{
			name: "Points far apart",
			x1:   -5, y1: -5,
			x2: 5, y2: 5,
			radius: 1.0,
			want:   math.Exp(-100),
		},
		{
			name: "Large radius",
			x1:   0, y1: 0,
			x2: 10, y2: 10,
			radius: 100.0,
			want:   math.Exp(-0.01),
		},
		{
			name: "Very small radius",
			x1:   0, y1: 0,
			x2: 1, y2: 1,
			radius: 0.1,
			want:   math.Exp(-100),
		},
	}

	for _, tt := range tests {
		t.Run("Gaussian: "+tt.name, func(t *testing.T) {
			got := g.Weight(tt.x1, tt.y1, tt.x2, tt.y2, tt.radius)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Weight(%d, %d, %d, %d, %f) = %v, want %v", tt.x1, tt.y1, tt.x2, tt.y2, tt.radius, got, tt.want)
			}
		})

		t.Run("CutGaussian: "+tt.name, func(t *testing.T) {
			got := g2.Weight(tt.x1, tt.y1, tt.x2, tt.y2, tt.radius)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Weight(%d, %d, %d, %d, %f) = %v, want %v", tt.x1, tt.y1, tt.x2, tt.y2, tt.radius, got, tt.want)
			}
		})
	}
}
func TestBoxWeight(t *testing.T) {
	b := &neighborhood.Box{}

	tests := []struct {
		name   string
		x1, y1 int
		x2, y2 int
		radius float64
		want   float64
	}{
		{
			name: "Inside box",
			x1:   0, y1: 0,
			x2: 2, y2: 2,
			radius: 3.0,
			want:   1.0,
		},
		{
			name: "On box edge",
			x1:   0, y1: 0,
			x2: 3, y2: 0,
			radius: 3.0,
			want:   1.0,
		},
		{
			name: "Outside box",
			x1:   0, y1: 0,
			x2: 4, y2: 0,
			radius: 3.0,
			want:   0.0,
		},
		{
			name: "Diagonal inside",
			x1:   0, y1: 0,
			x2: 2, y2: 2,
			radius: 2.9,
			want:   1.0,
		},
		{
			name: "Diagonal outside",
			x1:   0, y1: 0,
			x2: 3, y2: 3,
			radius: 4.0,
			want:   0.0,
		},
		{
			name: "Zero radius",
			x1:   0, y1: 0,
			x2: 0, y2: 0,
			radius: 0.0,
			want:   1.0,
		},
		{
			name: "Negative coordinates",
			x1:   -2, y1: -2,
			x2: -3, y2: -3,
			radius: 3.0,
			want:   1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.Weight(tt.x1, tt.y1, tt.x2, tt.y2, tt.radius)
			if got != tt.want {
				t.Errorf("Weight(%d, %d, %d, %d, %f) = %v, want %v in %s", tt.x1, tt.y1, tt.x2, tt.y2, tt.radius, got, tt.want, tt.name)
			}
		})
	}
}
