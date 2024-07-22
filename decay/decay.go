package decay

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var decays = map[string]func() Decay{}

func init() {
	d := []func() Decay{
		func() Decay { return &Constant{} },
		func() Decay { return &Linear{} },
		func() Decay { return &Power{} },
		func() Decay { return &Polynomial{} },
	}
	for _, v := range d {
		vv := v()
		if _, ok := decays[vv.Name()]; ok {
			panic("duplicate decay name: " + vv.Name())
		}
		decays[vv.Name()] = v
	}
}

func FromString(nameAndArgs string) (Decay, error) {
	parts := strings.Split(nameAndArgs, " ")

	dFunc, ok := decays[parts[0]]
	if !ok {
		return nil, fmt.Errorf("unknown decay: %s", parts[0])
	}
	d := dFunc()
	if len(parts) == 1 {
		return d, nil
	}
	args := make([]float64, len(parts)-1)
	for i := 1; i < len(parts); i++ {
		v, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			return nil, err
		}
		args[i-1] = v
	}
	err := d.SetArgs(args...)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type Decay interface {
	Name() string
	Decay(epoch, total int) float64
	SetArgs(args ...float64) error
}

type Constant struct {
	Value float64
}

func (c *Constant) Name() string {
	return "constant"
}

func (c *Constant) Decay(epoch, total int) float64 {
	return c.Value
}

func (c *Constant) SetArgs(args ...float64) error {
	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg, got %d", len(args))
	}
	c.Value = args[0]
	return nil
}

type Linear struct {
	Start float64
	End   float64
}

func (l *Linear) Name() string {
	return "linear"
}

func (l *Linear) Decay(epoch, total int) float64 {
	d := 1 - float64(epoch)/float64(total)
	return l.End + d*(l.Start-l.End)
}

func (l *Linear) SetArgs(args ...float64) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 args, got %d", len(args))
	}
	l.Start = args[0]
	l.End = args[1]
	return nil
}

type Power struct {
	Start float64
	End   float64
}

func (p *Power) Name() string {
	return "power"
}

func (p *Power) Decay(epoch, total int) float64 {
	d := float64(epoch) / float64(total)
	return p.Start * math.Pow(p.End/p.Start, d)
}

func (p *Power) SetArgs(args ...float64) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 args, got %d", len(args))
	}
	p.Start = args[0]
	p.End = args[1]
	return nil
}

type Polynomial struct {
	Start float64
	End   float64
	Exp   float64
}

func (p *Polynomial) Name() string {
	return "polynomial"
}

func (p *Polynomial) Decay(epoch, total int) float64 {
	d := float64(epoch) / float64(total)
	return p.End + (p.Start-p.End)*math.Pow(1-d, p.Exp)
}

func (p *Polynomial) SetArgs(args ...float64) error {
	if len(args) != 3 {
		return fmt.Errorf("expected 3 args, got %d", len(args))
	}
	p.Start = args[0]
	p.End = args[1]
	p.Exp = args[2]
	return nil
}
