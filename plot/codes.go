package plot

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/benoitmasson/plotters/piechart"
	"github.com/mlange-42/som"
	"github.com/mlange-42/som/norm"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type CodePlot interface {
	Plot(data []float64, dRange Range) (*plot.Plot, error)
}

func Codes(s *som.Som, columns [][2]int, normalized bool, zeroAxis bool, plotType CodePlot, size image.Point) (image.Image, error) {
	legendHeight := 20
	hPad, vPad := 2, 4

	codeHeight := (size.Y - legendHeight) / s.Size().Height
	codeWidth := size.X / s.Size().Width

	img := vgimg.NewWith(vgimg.UseWH(font.Length(size.X), font.Length(size.Y)), vgimg.UseDPI(72))
	dc := draw.New(img)

	dRange := dataRange(s, columns, normalized, zeroAxis)

	w, h := s.Size().Width, s.Size().Height
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := draw.Crop(dc,
				font.Length(x*codeWidth+hPad), font.Length((x+1-w)*codeWidth-hPad),
				font.Length(y*codeHeight+legendHeight+vPad), font.Length((y+1-h)*codeHeight-vPad))

			node := s.Size().Index(x, y)
			data := nodeData(s, node, columns, normalized)
			p, err := plotType.Plot(data, dRange)
			if err != nil {
				return nil, err
			}

			p.Draw(c)
		}
	}

	return img.Image(), nil
}

type CodeLines struct{}

func (c *CodeLines) Plot(data []float64, dataRange Range) (*plot.Plot, error) {
	p := plot.New()

	lines, err := plotter.NewLine(SimpleXY(data))
	if err != nil {
		return nil, err
	}

	cleanupAxes(p)
	p.Y.AutoRescale = false
	p.Y.Min, p.Y.Max = dataRange.Min, dataRange.Max

	p.Add(lines)

	return p, nil
}

type CodePie struct{}

func (c *CodePie) Plot(data []float64, dataRange Range) (*plot.Plot, error) {
	for _, v := range data {
		if v < 0 {
			return nil, fmt.Errorf("negative values not supported in pie chart")
		}
	}

	p := plot.New()
	p.HideAxes()

	pie, err := piechart.NewPieChart(plotter.Values(data))
	if err != nil {
		return nil, err
	}
	pie.Labels.Show = false
	pie.Radius = 1
	pie.LineStyle.Width = 1
	pie.LineStyle.Color = color.White
	pie.Color = color.Black

	p.Add(pie)

	return p, nil
}

func cleanupAxes(p *plot.Plot) {
	p.X.Tick.Marker = plot.TickerFunc(func(min float64, max float64) []plot.Tick { return nil })
	p.Y.Tick.Label.Font.Size = 0
	p.Y.Tick.Length = 2

	p.X.Padding = 0
	p.Y.Padding = 0
}

func dataRange(s *som.Som, columns [][2]int, normalized bool, zeroAxis bool) Range {
	minValue := math.Inf(1)
	maxValue := math.Inf(-1)

	identity := norm.Identity{}

	nodes := s.Size().Nodes()
	for _, c := range columns {
		lay := s.Layers()[c[0]]
		var normalizer norm.Normalizer
		if normalized {
			normalizer = &identity
		} else {
			normalizer = lay.Normalizers()[c[1]]
		}
		for i := 0; i < nodes; i++ {
			value := normalizer.DeNormalize(lay.GetAt(i, c[1]))
			minValue = math.Min(minValue, value)
			maxValue = math.Max(maxValue, value)
		}
	}

	if zeroAxis {
		minValue = math.Min(minValue, 0)
		maxValue = math.Max(maxValue, 0)
	}

	return Range{minValue, maxValue}
}

func nodeData(s *som.Som, node int, columns [][2]int, normalized bool) []float64 {
	data := make([]float64, len(columns))

	if normalized {
		for i, c := range columns {
			lay := s.Layers()[c[0]]
			data[i] = lay.GetAt(node, c[1])
		}
	} else {
		for i, c := range columns {
			lay := s.Layers()[c[0]]
			norm := lay.Normalizers()[c[1]]
			data[i] = norm.DeNormalize(lay.GetAt(node, c[1]))
		}
	}

	return data
}

type Range struct {
	Min, Max float64
}
