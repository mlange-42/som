package plot

import (
	"fmt"
	"image"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

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
	legendWidth := l.TextStyle.Rectangle("9999.00").Max.X + l.ThumbnailWidth + 3*vg.Millimeter
	l.YOffs = font.Length(-titleHeight) // Adjust the legend down a little.
	l.XOffs = -2 * vg.Millimeter

	dcPlot := draw.Crop(dc, 0, -legendWidth-vg.Millimeter, 0, 0) // Make space for the legend.
	p.Draw(dcPlot)
	l.Draw(dc)

	return img.Image(), nil
}

func createLabels(labels []string, positions []plotter.XY, baseStyle text.Style) *ZeroSizeLabel {
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
	return &ZeroSizeLabel{l}
}
