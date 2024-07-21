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

// Neighborhood is an interface that defines the behavior of a neighborhood function.
// The Name method returns the name of the neighborhood.
// The Weight method returns the weight of a point at the given distance from the center, based on the given radius.
// The MaxRadius method returns the maximum radius for which the neighborhood function is non-zero.
type Neighborhood interface {
	Name() string
	Weight(distance, radius float64) float64
	MaxRadius(radius float64) int
}

// Gaussian implements [Neighborhood] for the Gaussian neighborhood function.
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

// CutGaussian implements [Neighborhood] for the cut Gaussian neighborhood function.
// It returns 0 if the distance is greater than the radius (i.e. SD of the Gaussian kernel).
type CutGaussian struct{}

func (g *CutGaussian) Name() string {
	return "cutgaussian"
}

func (g *CutGaussian) Weight(distance, radius float64) float64 {
	if distance > radius {
		return 0
	}
	return math.Exp(-(distance * distance) / (2 * radius * radius))
}

func (g *CutGaussian) MaxRadius(radius float64) int {
	return int(radius)
}

// Linear implements [Neighborhood] for the linear neighborhood function.
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

// YourNeighborhood implements [Neighborhood] for a box or constant neighborhood function.
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
