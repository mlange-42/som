package plot

import (
	"image/color"

	somplotter "github.com/mlange-42/som/plot/plotter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type CodeBar struct {
	Colors     []color.Color
	Horizontal bool
	AdjustAxis bool
}

func (c *CodeBar) Plot(data []float64, dataRange Range) (*plot.Plot, []plot.Thumbnailer, error) {
	if len(c.Colors) == 0 {
		c.Colors = DefaultColors
	}

	p := plot.New()

	cleanupAxes(p, c.Horizontal)
	if c.Horizontal {
		p.HideY()
	} else {
		p.HideX()
	}

	if c.AdjustAxis {
		ax := &p.Y
		if c.Horizontal {
			ax = &p.X
		}
		ax.AutoRescale = false
		ax.Min, ax.Max = dataRange.Min, dataRange.Max
	}

	thumbs := make([]plot.Thumbnailer, len(data))
	for i, v := range data {
		bar, err := somplotter.NewBarChart(plotter.Values([]float64{v}), 1)
		if err != nil {
			return nil, nil, err
		}
		bar.Horizontal = c.Horizontal
		bar.Color = c.Colors[i%len(c.Colors)]
		bar.XMin = float64(i)

		p.Add(bar)

		thumbs[i] = bar
	}

	return p, thumbs, nil
}
