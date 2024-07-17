package norm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mlange-42/som/table"
)

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

func FromString(nameAndArgs string) (Normalizer, error) {
	parts := strings.Split(nameAndArgs, " ")

	nFunc, ok := normalizers[parts[0]]
	if !ok {
		return nil, fmt.Errorf("unknown normalizer: %s", parts[0])
	}
	n := nFunc()
	if len(parts) == 1 {
		return n, nil
	}
	args := make([]float64, len(parts)-1)
	for i := 1; i < len(parts); i++ {
		v, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			return nil, err
		}
		args[i-1] = v
	}
	err := n.SetArgs(args...)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func ToString(n Normalizer) string {
	args := n.GetArgs()
	if len(args) == 0 {
		return n.Name()
	}
	s := n.Name() + " "
	for i, v := range args {
		s += strconv.FormatFloat(v, 'f', -1, 64)
		if i < len(args)-1 {
			s += " "
		}
	}
	return s
}

type Normalizer interface {
	Name() string
	Normalize(value float64) float64
	DeNormalize(value float64) float64
	Initialize(t *table.Table, column int)
	SetArgs(args ...float64) error
	GetArgs() []float64
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

func (n *None) SetArgs(args ...float64) error {
	return nil
}

func (n *None) GetArgs() []float64 {
	return nil
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

func (g *Gaussian) SetArgs(args ...float64) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", len(args))
	}
	g.mean = args[0]
	g.std = args[1]
	return nil
}

func (g *Gaussian) GetArgs() []float64 {
	return []float64{g.mean, g.std}
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

func (u *Uniform) SetArgs(args ...float64) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", len(args))
	}
	u.min = args[0]
	u.max = args[1]
	return nil
}

func (u *Uniform) GetArgs() []float64 {
	return []float64{u.min, u.max}
}
