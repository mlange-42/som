package plot

import (
	"fmt"
	"image"

	"github.com/mlange-42/som"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
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

func Heatmap(som *som.Som, layer, column, width, height int) (image.Image, error) {
	g := grid{
		som:    som,
		layer:  layer,
		column: column,
	}

	colors := 12
	//pal := palette.Heat(colors, 1)
	pal := palette.Rainbow(colors, palette.Blue, palette.Red, 1, 1, 1)
	h := plotter.NewHeatMap(&g, pal)

	p := plot.New()
	p.Title.Text = fmt.Sprintf("%s - %s", som.Layers()[layer].Name(), som.Layers()[layer].ColumnNames()[column])
	p.HideAxes()

	p.Add(h)

	l := plot.NewLegend()
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
		l.Add(fmt.Sprintf("%.2g", val), t)
	}

	img := vgimg.New(font.Length(width), font.Length(height))
	dc := draw.New(img)

	l.Top = true
	// Calculate the width of the legend.
	r := l.Rectangle(dc)
	legendWidth := r.Max.X - r.Min.X
	l.YOffs = -p.Title.TextStyle.FontExtents().Height // Adjust the legend down a little.

	l.Draw(dc)
	dc = draw.Crop(dc, 0, -legendWidth-vg.Millimeter, 0, 0) // Make space for the legend.
	p.Draw(dc)

	return img.Image(), nil

	/*w, err := os.Create("heatMap.png")
	if err != nil {
		return err
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		return err
	}

	return nil*/
}
