package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type CodeLines struct{}

func (c *CodeLines) Plot(data []float64, dataRange Range) (*plot.Plot, []plot.Thumbnailer, error) {
	p := plot.New()

	lines, err := plotter.NewLine(SimpleXY(data))
	if err != nil {
		return nil, nil, err
	}

	cleanupAxes(p)
	p.Y.AutoRescale = false
	p.Y.Min, p.Y.Max = dataRange.Min, dataRange.Max

	p.Add(lines)

	return p, nil, nil
}
