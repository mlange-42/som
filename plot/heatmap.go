package plot

import (
	"fmt"
	"image"

	"github.com/mlange-42/som"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type grid struct {
	som    *som.Som
	layer  int
	column int
}

func (g *grid) Dims() (c, r int) {
	return g.som.Size().Width, g.som.Size().Height
}

func (g *grid) Z(c, r int) float64 {
	l := &g.som.Layers()[g.layer]
	v := l.Get(c, r, g.column)
	return l.Normalizers()[g.column].DeNormalize(v)
}

func (g *grid) X(c int) float64 {
	return float64(c)
}

func (g *grid) Y(r int) float64 {
	return float64(r)
}

func Heatmap(som *som.Som, layer, column, width, height int, labels []string, positions []plotter.XY) (image.Image, error) {
	g := grid{
		som:    som,
		layer:  layer,
		column: column,
	}
	p := plot.New()
	l := plot.NewLegend()
	p.Title.TextStyle.Font.Size = 16

	titleHeight := p.Title.TextStyle.FontExtents().Height.Points()
	legendHeight := height - int(titleHeight) - 2
	itemHeight := l.TextStyle.Rectangle("Aq").Max.Y.Points() + 1
	numColors := legendHeight / int(itemHeight)

	pal := palette.Rainbow(numColors, palette.Blue, palette.Red, 1, 1, 1)
	h := plotter.NewHeatMap(&g, pal)

	p.Title.Text = fmt.Sprintf("%s - %s", som.Layers()[layer].Name(), som.Layers()[layer].ColumnNames()[column])
	p.HideAxes()
	p.Add(h)

	labelsPlot := createLabels(labels, positions, l.TextStyle)
	p.Add(labelsPlot)

	thumbs := plotter.PaletteThumbnailers(pal)
	for i := len(thumbs) - 1; i >= 0; i-- {
		t := thumbs[i]
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
