package norm

import "github.com/mlange-42/som/table"

var normalizers = map[string]Normalizer{}

func init() {
	n := []Normalizer{
		&Gaussian{},
		&Uniform{},
	}
	for _, v := range n {
		if _, ok := normalizers[v.Name()]; ok {
			panic("duplicate normalizer name: " + v.Name())
		}
		normalizers[v.Name()] = v
	}
}

func GetNormalizer(name string) (Normalizer, bool) {
	n, ok := normalizers[name]
	return n, ok
}

type Normalizer interface {
	Name() string
	Normalize(value float64) float64
	DeNormalize(value float64) float64
	Initialize(t *table.Table, column int)
}

type Gaussian struct {
	mean, std float64
}

func (g *Gaussian) Name() string {
	return "gaussian"
}

func (g *Gaussian) Normalize(value float64) float64 {
	return (value - g.mean) / g.std
}

func (g *Gaussian) DeNormalize(value float64) float64 {
	return value*g.std + g.mean
}

func (g *Gaussian) Initialize(t *table.Table, column int) {
	g.mean = t.Mean(column)
	g.std = t.StdDev(column)
}

type Uniform struct {
	min, max float64
}

func (u *Uniform) Name() string {
	return "uniform"
}

func (u *Uniform) Normalize(value float64) float64 {
	return (value - u.min) / (u.max - u.min)
}

func (u *Uniform) DeNormalize(value float64) float64 {
	return value*(u.max-u.min) + u.min
}

func (u *Uniform) Initialize(t *table.Table, column int) {
	u.min, u.max = t.Range(column)
}
