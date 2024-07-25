package plot

import (
	"image"
	"math"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/norm"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type CodePlot interface {
	Plot(data []float64, dRange Range) (*plot.Plot, []plot.Thumbnailer, error)
}

func Codes(s *som.Som, columns [][2]int,
	boundariesLayer int,
	normalized bool, zeroAxis bool,
	plotType CodePlot, size image.Point) (image.Image, error) {
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
		l.AdjustColumns(font.Length(size.X - 12))
		l.XOffs = (font.Length(size.X) - l.Rectangle(dc).Max.X) / 2
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

	if boundariesLayer >= 0 {
		p := plot.New()
		p.BackgroundColor = image.Transparent
		p.HideAxes()

		c := draw.Crop(dc, 0, 0, font.Length(legendHeight), 0)

		_, classIndices := conv.LayerToClasses(s.Layers()[boundariesLayer])
		bounds := &IntGrid{Size: *s.Size(), Values: classIndices}
		bound, err := NewGridBoundaries(bounds)
		if err != nil {
			return nil, err
		}
		p.Add(bound)

		p.Draw(c)
	}

	return img.Image(), nil
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
