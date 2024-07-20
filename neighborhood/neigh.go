package neighborhood

import "math"

var neighborhoods = map[string]Neighborhood{}

func init() {
	n := []Neighborhood{
		&Gaussian{},
		&CutGaussian{},
		&Linear{},
		&Box{},
	}
	for _, v := range n {
		if _, ok := neighborhoods[v.Name()]; ok {
			panic("duplicate neighborhood name: " + v.Name())
		}
		neighborhoods[v.Name()] = v
	}
}

func GetNeighborhood(name string) (Neighborhood, bool) {
	n, ok := neighborhoods[name]
	return n, ok
}

type Neighborhood interface {
	Name() string
	Weight(distance, radius float64) float64
	MaxRadius(radius float64) int
}

type Gaussian struct{}

func (g *Gaussian) Name() string {
	return "gaussian"
}

func (g *Gaussian) Weight(distance, radius float64) float64 {
	return math.Exp(-(distance * distance) / (2 * radius * radius))
}

func (g *Gaussian) MaxRadius(radius float64) int {
	return int(3 * radius)
}

type CutGaussian struct{}

func (g *CutGaussian) Name() string {
	return "cutgaussian"
}

func (g *CutGaussian) Weight(distance, radius float64) float64 {
	return math.Exp(-(distance * distance) / (2 * radius * radius))
}

func (g *CutGaussian) MaxRadius(radius float64) int {
	return int(radius)
}

type Linear struct{}

func (l *Linear) Name() string {
	return "linear"
}

func (l *Linear) Weight(distance, radius float64) float64 {
	if distance <= radius {
		return 1 - distance/radius
	}
	return 0
}

func (g *Linear) MaxRadius(radius float64) int {
	return int(radius)
}

type Box struct{}

func (b *Box) Name() string {
	return "box"
}

func (b *Box) Weight(distance, radius float64) float64 {
	if distance <= radius {
		return 1
	}
	return 0
}

func (b *Box) MaxRadius(radius float64) int {
	return int(radius)
}
