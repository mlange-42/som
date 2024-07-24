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
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type CodePlot interface {
	Plot(data []float64, dRange Range) (*plot.Plot, []plot.Thumbnailer, error)
}

func Codes(s *som.Som, columns [][2]int, normalized bool, zeroAxis bool, plotType CodePlot, size image.Point) (image.Image, error) {
	legendFontSize := 16

	img := vgimg.NewWith(vgimg.UseWH(font.Length(size.X), font.Length(size.Y)), vgimg.UseDPI(72))
	dc := draw.New(img)

	plots := make([]*plot.Plot, s.Size().Width*s.Size().Height)

	dRange := dataRange(s, columns, normalized, zeroAxis)
	w, h := s.Size().Width, s.Size().Height

	var thumbs []plot.Thumbnailer
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			node := s.Size().Index(x, y)
			data := nodeData(s, node, columns, normalized)

			var p *plot.Plot
			var err error
			p, thumbs, err = plotType.Plot(data, dRange)
			if err != nil {
				return nil, err
			}
			plots[node] = p
		}
	}

	var l Legend
	if len(thumbs) > 0 {
		l = NewLegend()
		l.Left = true
		l.YOffs = vg.Millimeter * 2
		l.TextStyle.Font.Size = font.Length(legendFontSize)
		for i := range thumbs {
			c := columns[i]
			label := s.Layers()[c[0]].ColumnNames()[c[1]]
			l.Add(label, thumbs[i])
		}
		l.AdjustColumns(font.Length(size.X))
		l.XOffs = (font.Length(size.X) - l.Rectangle(dc).Max.X) / 2

		l.Draw(dc)
	}

	legendHeight := (legendFontSize + 2) * int(math.Ceil(float64(len(thumbs))/float64(l.Columns)))
	hPad, vPad := 2, 4

	codeHeight := (size.Y - legendHeight) / s.Size().Height
	codeWidth := size.X / s.Size().Width

	for i, p := range plots {
		x, y := s.Size().Coords(i)
		c := draw.Crop(dc,
			font.Length(x*codeWidth+hPad), font.Length((x+1-w)*codeWidth-hPad),
			font.Length(y*codeHeight+legendHeight+vPad), font.Length((y+1-h)*codeHeight-vPad))

		p.Draw(c)
	}

	if len(thumbs) > 0 {
		l.Draw(dc)
	}

	return img.Image(), nil
}

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
