package decay

import "math"

type Decay interface {
	Decay(epoch, total int) float64
}

type Linear struct {
	Start float64
	End   float64
}

func (l *Linear) Decay(epoch, total int) float64 {
	d := 1 - float64(epoch)/float64(total)
	return l.End + d*(l.Start-l.End)
}

type Power struct {
	Start float64
	End   float64
}

func (p *Power) Decay(epoch, total int) float64 {
	d := float64(epoch) / float64(total)
	return p.Start * math.Pow(p.End/p.Start, d)
}
