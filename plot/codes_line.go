package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type CodeLines struct {
	StepStyle plotter.StepKind
	Vertical  bool
}

func (c *CodeLines) Plot(data []float64, dataRange Range) (*plot.Plot, []plot.Thumbnailer, error) {
	p := plot.New()

	xy := SimpleXY{
		Values:   data,
		Vertical: c.Vertical,
	}
	lines, err := plotter.NewLine(&xy)
	if err != nil {
		return nil, nil, err
	}
	lines.StepStyle = c.StepStyle

	cleanupAxes(p)
	p.Y.AutoRescale = false
	p.Y.Min, p.Y.Max = dataRange.Min, dataRange.Max

	p.Add(lines)

	return p, nil, nil
}
