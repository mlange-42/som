package plot

import (
	"image"
	"image/color"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/layer"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func XY(
	title string, g plotter.XYer, size layer.Size, width, height int,
	categories []string, catIndices []int, drawGrid bool,
	labels []string, positions []plotter.XY,
) (image.Image, error) {

	p := plot.New()
	l := plot.NewLegend()
	p.Title.TextStyle.Font.Size = 16

	titleHeight := p.Title.TextStyle.FontExtents().Height.Points()

	pal := NewRandomPalette(len(categories))
	h, err := plotter.NewScatter(g)
	if err != nil {
		return nil, err
	}
	if len(categories) > 0 {
		h.GlyphStyleFunc = func(i int) draw.GlyphStyle {
			cat := catIndices[i]
			return draw.GlyphStyle{
				Color:  pal.Colors()[cat],
				Shape:  draw.CircleGlyph{},
				Radius: vg.Length(3),
			}
		}
	} else {
		h.GlyphStyle = draw.GlyphStyle{
			Shape:  draw.CircleGlyph{},
			Radius: vg.Length(3),
		}
	}

	if drawGrid {
		err := addMapGrid(size, p, g)
		if err != nil {
			return nil, err
		}
	}

	p.Title.Text = title
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

func addMapGrid(size layer.Size, p *plot.Plot, g plotter.XYer) error {
	ls := draw.LineStyle{
		Color: color.NRGBA{R: 160, G: 160, B: 160, A: 255},
		Width: vg.Length(0.5),
	}

	for row := 0; row < size.Height; row++ {
		xy := SomRowColXY{
			Xy:    g,
			Size:  size,
			Index: row,
			IsRow: true,
		}
		line, err := plotter.NewLine(&xy)
		if err != nil {
			return err
		}
		line.LineStyle = ls
		p.Add(line)
	}

	for col := 0; col < size.Width; col++ {
		xy := SomRowColXY{
			Xy:    g,
			Size:  size,
			Index: col,
			IsRow: false,
		}
		line, err := plotter.NewLine(&xy)
		if err != nil {
			return err
		}
		line.LineStyle = ls
		p.Add(line)
	}

	return nil
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

type SomRowColXY struct {
	Xy    plotter.XYer
	Size  layer.Size
	Index int
	IsRow bool
}

func (s *SomRowColXY) XY(i int) (x, y float64) {
	var col, row int
	if s.IsRow {
		col = i
		row = s.Index
	} else {
		col = s.Index
		row = i
	}
	idx := s.Size.Index(col, row)
	return s.Xy.XY(idx)
}

func (s *SomRowColXY) Len() int {
	if s.IsRow {
		return s.Size.Width
	}
	return s.Size.Height
}
