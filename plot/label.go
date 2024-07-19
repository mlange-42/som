package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

type ZeroSizeLabel struct {
	labels *plotter.Labels
}

func (l *ZeroSizeLabel) Plot(c draw.Canvas, p *plot.Plot) {
	l.labels.Plot(c, p)
}
