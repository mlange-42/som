package plot

import (
	"image/color"

	"github.com/benoitmasson/plotters/piechart"
	"gonum.org/v1/plot"
)

type CodeRose struct {
	Colors []color.Color
}

func (c *CodeRose) Plot(data []float64, dataRange Range) (*plot.Plot, []plot.Thumbnailer, error) {
	if len(c.Colors) == 0 {
		c.Colors = DefaultColors
	}

	p := plot.New()
	p.HideAxes()

	thumbs := make([]plot.Thumbnailer, len(data))

	for i, v := range data {
		pie, err := piechart.NewPieChart(&ConstantValues{Val: 1, Length: 1})
		if err != nil {
			return nil, nil, err
		}
		pie.Radius = (v - dataRange.Min) / (dataRange.Max - dataRange.Min)

		pie.Labels.Show = false
		pie.LineStyle.Width = 0.5
		pie.LineStyle.Color = color.White
		pie.Color = c.Colors[i%len(c.Colors)]

		pie.Total = float64(len(data))
		pie.Offset.Value = float64(i)
		p.Add(pie)

		thumbs[i] = pie
	}

	return p, thumbs, nil
}
