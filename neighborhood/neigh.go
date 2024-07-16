package neighborhood

import "math"

type Neighborhood interface {
	Weight(x1, y1, x2, y2 int, radius float64) float64
	MaxRadius(radius float64) int
}

type Gaussian struct{}

func (g *Gaussian) Weight(x1, y1, x2, y2 int, radius float64) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Exp(-(dx*dx + dy*dy) / (2 * radius * radius))
}

func (g *Gaussian) MaxRadius(radius float64) int {
	return -1
}

type CutGaussian struct{}

func (g *CutGaussian) Weight(x1, y1, x2, y2 int, radius float64) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Exp(-(dx*dx + dy*dy) / (2 * radius * radius))
}

func (g *CutGaussian) MaxRadius(radius float64) int {
	return int(radius)
}

type Linear struct{}

func (l *Linear) Weight(x1, y1, x2, y2 int, radius float64) float64 {
	dx, dy := float64(x1-x2), float64(y1-y2)
	distSq := dx*dx + dy*dy
	if distSq <= radius*radius {
		return 1 - math.Sqrt(distSq)/radius
	}
	return 0
}

func (g *Linear) MaxRadius(radius float64) int {
	return int(radius)
}

type Box struct{}

func (b *Box) Weight(x1, y1, x2, y2 int, radius float64) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	if dx*dx+dy*dy <= radius*radius {
		return 1
	}
	return 0
}

func (b *Box) MaxRadius(radius float64) int {
	return int(radius)
}
