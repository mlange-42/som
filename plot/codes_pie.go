package plot

import (
	"fmt"
	"image/color"

	"github.com/benoitmasson/plotters/piechart"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type CodePie struct {
	Colors []color.Color
}

func (c *CodePie) Plot(data []float64, dataRange Range) (*plot.Plot, []plot.Thumbnailer, error) {
	if len(c.Colors) == 0 {
		c.Colors = DefaultColors
	}

	total := 0.0
	for _, v := range data {
		if v < 0 {
			return nil, nil, fmt.Errorf("negative values not supported in pie chart")
		}
		total += v
	}

	p := plot.New()
	p.HideAxes()

	thumbs := make([]plot.Thumbnailer, len(data))

	offset := 0.0
	for i, v := range data {
		pie, err := piechart.NewPieChart(plotter.Values([]float64{v}))
		if err != nil {
			return nil, nil, err
		}
		pie.Labels.Show = false
		pie.Radius = 1
		pie.LineStyle.Width = 1
		pie.LineStyle.Color = color.White
		pie.Color = c.Colors[i%len(c.Colors)]

		pie.Total = total
		pie.Offset.Value = offset
		p.Add(pie)

		thumbs[i] = pie

		offset += v
	}

	return p, thumbs, nil
}
