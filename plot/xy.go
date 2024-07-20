package plot

import (
	"image"

	"github.com/mlange-42/som"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func XY(title string, g plotter.XYer, width, height int, categories []string, labels []string, positions []plotter.XY) (image.Image, error) {
	p := plot.New()
	l := plot.NewLegend()
	p.Title.TextStyle.Font.Size = 16

	titleHeight := p.Title.TextStyle.FontExtents().Height.Points()

	pal := NewRandomPalette(len(categories))
	h, err := plotter.NewScatter(g)
	if err != nil {
		return nil, err
	}

	p.Title.Text = title
	//p.HideAxes()
	p.Add(h)

	thumbs := plotter.PaletteThumbnailers(pal)
	for i := len(thumbs) - 1; i >= 0; i-- {
		t := thumbs[i]
		l.Add(categories[i], t)
	}

	img := vgimg.NewWith(vgimg.UseWH(font.Length(width), font.Length(height)), vgimg.UseDPI(72))
	dc := draw.New(img)

	l.Top = true
	legendWidth := l.TextStyle.Rectangle("9999.00").Max.X + l.ThumbnailWidth + 3*vg.Millimeter
	l.YOffs = font.Length(-titleHeight) // Adjust the legend down a little.
	l.XOffs = -2 * vg.Millimeter

	l.Draw(dc)

	dc = draw.Crop(dc, 0, -legendWidth-vg.Millimeter, 0, 0) // Make space for the legend.
	p.Draw(dc)

	return img.Image(), nil
}

type SomXY struct {
	Som     *som.Som
	XLayer  int
	XColumn int
	YLayer  int
	YColumn int
}

func (s *SomXY) XY(i int) (x, y float64) {
	return s.Som.Layers()[s.XLayer].GetAt(i, s.XColumn), s.Som.Layers()[s.YLayer].GetAt(i, s.YColumn)
}

func (s *SomXY) Len() int {
	return s.Som.Size().Nodes()
}
