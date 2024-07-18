package plot

import (
	"fmt"
	"image"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/layer"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type SomLayerGrid struct {
	Som    *som.Som
	Layer  int
	Column int
}

func (g *SomLayerGrid) Dims() (c, r int) {
	return g.Som.Size().Width, g.Som.Size().Height
}

func (g *SomLayerGrid) Z(c, r int) float64 {
	l := &g.Som.Layers()[g.Layer]
	v := l.Get(c, r, g.Column)
	return l.Normalizers()[g.Column].DeNormalize(v)
}

func (g *SomLayerGrid) X(c int) float64 {
	return float64(c)
}

func (g *SomLayerGrid) Y(r int) float64 {
	return float64(r)
}

type ClassesGrid struct {
	Size    layer.Size
	Indices []int
}

func (g *ClassesGrid) Dims() (c, r int) {
	return g.Size.Width, g.Size.Height
}

func (g *ClassesGrid) Z(c, r int) float64 {
	idx := r + c*g.Size.Height
	return float64(g.Indices[idx])
}

func (g *ClassesGrid) X(c int) float64 {
	return float64(c)
}

func (g *ClassesGrid) Y(r int) float64 {
	return float64(r)
}

func Heatmap(title string, g plotter.GridXYZ, width, height int, categories []string, labels []string, positions []plotter.XY) (image.Image, error) {
	p := plot.New()
	l := plot.NewLegend()
	p.Title.TextStyle.Font.Size = 16

	titleHeight := p.Title.TextStyle.FontExtents().Height.Points()
	legendHeight := height - int(titleHeight) - 2
	itemHeight := l.TextStyle.Rectangle("Aq").Max.Y.Points() + 1
	numColors := legendHeight / int(itemHeight)

	categorical := len(categories) > 0
	var pal palette.Palette
	if categorical {
		pal = NewRandomPalette(len(categories))
	} else {
		pal = palette.Rainbow(numColors, palette.Blue, palette.Red, 1, 1, 1)
	}
	h := plotter.NewHeatMap(g, pal)

	p.Title.Text = title
	p.HideAxes()
	p.Add(h)

	labelsPlot := createLabels(labels, positions, l.TextStyle)
	p.Add(labelsPlot)

	thumbs := plotter.PaletteThumbnailers(pal)
	for i := len(thumbs) - 1; i >= 0; i-- {
		t := thumbs[i]
		if categorical {
			l.Add(categories[i], t)
			continue
		}

		if i != 0 && i != len(thumbs)-1 {
			l.Add("", t)
			continue
		}
		var val float64
		switch i {
		case 0:
			val = h.Min
		case len(thumbs) - 1:
			val = h.Max
		}
		if val < 10000 {
			l.Add(fmt.Sprintf("%.2f", val), t)
		} else {
			l.Add(fmt.Sprintf("%.2g", val), t)
		}
	}

	img := vgimg.NewWith(vgimg.UseWH(font.Length(width), font.Length(height)), vgimg.UseDPI(72))
	dc := draw.New(img)

	l.Top = true
	// Calculate the width of the legend.
	//r := l.Rectangle(dc)
	legendWidth := l.TextStyle.Rectangle("9999.00").Max.X + l.ThumbnailWidth //r.Max.X - r.Min.X
	l.YOffs = font.Length(-titleHeight)                                      // Adjust the legend down a little.

	l.Draw(dc)
	dc = draw.Crop(dc, 0, -legendWidth-vg.Millimeter, 0, 0) // Make space for the legend.
	p.Draw(dc)

	return img.Image(), nil
}

func createLabels(labels []string, positions []plotter.XY, baseStyle text.Style) *Labels {
	style := baseStyle
	style.Font.Size = 12
	style.XAlign = text.XCenter
	style.YAlign = text.YCenter

	styles := make([]text.Style, len(labels))
	for i := range styles {
		styles[i] = style
	}

	l := &plotter.Labels{
		XYs:       positions,
		Labels:    labels,
		TextStyle: styles,
	}
	return &Labels{l}
}

type Labels struct {
	labels *plotter.Labels
}

func (l *Labels) Plot(c draw.Canvas, p *plot.Plot) {
	l.labels.Plot(c, p)
}
