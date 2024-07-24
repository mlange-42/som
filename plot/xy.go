package plot

import (
	"image"
	"image/color"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/norm"
	"github.com/mlange-42/som/table"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func XY(
	title string, g plotter.XYer, size layer.Size, width, height int,
	categories []string, catIndices []int, drawGrid bool,
	data plotter.XYer, dataCats []string, dataIndices []int,
	dataLegend bool,
) (image.Image, error) {

	p := plot.New()
	l := plot.NewLegend()
	p.Title.TextStyle.Font.Size = 16

	titleHeight := p.Title.TextStyle.FontExtents().Height.Points()

	var dataPalette palette.Palette
	if data != nil {
		dataScatter, err := plotter.NewScatter(data)
		if err != nil {
			return nil, err
		}
		dataPalette = setMarkerStyle(dataScatter, dataCats, dataIndices, 1.2, color.NRGBA{R: 80, G: 100, B: 240, A: 255})
		p.Add(dataScatter)
	}

	nodesScatter, err := plotter.NewScatter(g)
	if err != nil {
		return nil, err
	}
	nodePalette := setMarkerStyle(nodesScatter, categories, catIndices, 2.5, color.Black)

	if drawGrid {
		err := addMapGrid(size, p, g)
		if err != nil {
			return nil, err
		}
	}

	p.Title.Text = title
	p.Add(nodesScatter)

	if len(categories) > 0 {
		thumbs := plotter.PaletteThumbnailers(nodePalette)
		for i := len(thumbs) - 1; i >= 0; i-- {
			t := thumbs[i]
			l.Add(categories[i], t)
		}
	}
	if dataPalette != nil && (dataLegend || len(categories) == 0) {
		filler := plotter.PaletteThumbnailers(&WhitePalette{})
		l.Add("", filler[0])
		thumbs := plotter.PaletteThumbnailers(dataPalette)
		for i := len(thumbs) - 1; i >= 0; i-- {
			t := thumbs[i]
			l.Add(dataCats[i], t)
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

func setMarkerStyle(plot *plotter.Scatter, categories []string, catIndices []int, size float64, defaultColor color.Color) palette.Palette {
	pal := NewRandomPalette(DefaultColors, len(categories))
	if len(catIndices) > 0 {
		plot.GlyphStyleFunc = func(i int) draw.GlyphStyle {
			cat := catIndices[i]
			var col color.Color = color.RGBA{R: 100, G: 100, B: 100, A: 255}
			if cat >= 0 {
				col = pal.Colors()[cat]
			}
			return draw.GlyphStyle{
				Color:  col,
				Shape:  draw.CircleGlyph{},
				Radius: vg.Length(size),
			}
		}
	} else {
		plot.GlyphStyle = draw.GlyphStyle{
			Shape:  draw.CircleGlyph{},
			Radius: vg.Length(size),
			Color:  defaultColor,
		}
	}

	return pal
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
	lx := s.Som.Layers()[s.XLayer]
	ly := s.Som.Layers()[s.YLayer]
	vx, vy := lx.GetAt(i, s.XColumn), ly.GetAt(i, s.YColumn)
	return lx.Normalizers()[s.XColumn].DeNormalize(vx), ly.Normalizers()[s.YColumn].DeNormalize(vy)
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

type TableXY struct {
	XTable  *table.Table
	YTable  *table.Table
	XColumn int
	YColumn int
	XNorm   norm.Normalizer
	YNorm   norm.Normalizer
}

func (t *TableXY) XY(i int) (x, y float64) {
	vx, vy := t.XTable.Get(i, t.XColumn), t.YTable.Get(i, t.YColumn)
	return t.XNorm.DeNormalize(vx), t.YNorm.DeNormalize(vy)
}

func (t *TableXY) Len() int {
	return t.XTable.Rows()
}

type SimpleXY []float64

func (s SimpleXY) XY(i int) (x, y float64) {
	return float64(i), s[i]
}

func (s SimpleXY) Len() int {
	return len(s)
}
