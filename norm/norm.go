package norm

import "github.com/mlange-42/som/table"

var normalizers = map[string]func() Normalizer{}

func init() {
	n := []func() Normalizer{
		func() Normalizer { return &None{} },
		func() Normalizer { return &Gaussian{} },
		func() Normalizer { return &Uniform{} },
	}
	for _, v := range n {
		vv := v()
		if _, ok := normalizers[vv.Name()]; ok {
			panic("duplicate normalizer name: " + vv.Name())
		}
		normalizers[vv.Name()] = v
	}
}

func GetNormalizer(name string) (Normalizer, bool) {
	n, ok := normalizers[name]
	if ok {
		return n(), ok
	}
	return nil, ok
}

type Normalizer interface {
	Name() string
	Normalize(value float64) float64
	DeNormalize(value float64) float64
	Initialize(t *table.Table, column int)
}

type None struct{}

func (n *None) Name() string {
	return "none"
}

func (n *None) Normalize(value float64) float64 {
	return value
}

func (n *None) DeNormalize(value float64) float64 {
	return value
}

func (n *None) Initialize(t *table.Table, column int) {}

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
