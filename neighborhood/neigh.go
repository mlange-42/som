package neighborhood

import "math"

type Neighborhood interface {
	Weight(x1, y1, x2, y2 int, decay float64) float64
	Radius(decay float64) int
}

type Gaussian struct {
	Sigma float64
}

func (g *Gaussian) Weight(x1, y1, x2, y2 int, decay float64) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	s := g.Sigma * decay
	return math.Exp(-(dx*dx + dy*dy) / (2 * s * s))
}

func (g *Gaussian) Radius(decay float64) int {
	return -1
}

type CutGaussian struct {
	Sigma float64
}

func (g *CutGaussian) Weight(x1, y1, x2, y2 int, decay float64) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	s := g.Sigma * decay
	return math.Exp(-(dx*dx + dy*dy) / (2 * s * s))
}

func (g *CutGaussian) Radius(decay float64) int {
	return int(g.Sigma * decay)
}

type Box struct {
	Size float64
}

func (b *Box) Weight(x1, y1, x2, y2 int, decay float64) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	r := b.Size * decay
	if dx*dx+dy*dy <= r*r {
		return 1
	}
	return 0
}

func (b *Box) Radius(decay float64) int {
	return int(b.Size * decay)
}
