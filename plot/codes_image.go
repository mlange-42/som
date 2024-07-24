package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
)

type CodeImage struct {
	Rows int
}

func (c *CodeImage) Plot(data []float64, dataRange Range) (*plot.Plot, []plot.Thumbnailer, error) {
	p := plot.New()
	p.HideAxes()

	grid := NewImageGrid(data, c.Rows)
	pal := moreland.BlackBody().Palette(12)
	hm := plotter.NewHeatMap(&grid, pal)
	p.Add(hm)

	return p, nil, nil
}

type ImageGrid struct {
	data []float64
	cols int
	rows int
}

func NewImageGrid(data []float64, rows int) ImageGrid {
	return ImageGrid{
		data: data,
		cols: len(data) / rows,
		rows: rows,
	}
}

func (g *ImageGrid) Dims() (int, int) {
	return g.cols, g.rows
}

func (g *ImageGrid) X(c int) float64 {
	return float64(c)
}

func (g *ImageGrid) Y(r int) float64 {
	return float64(r)
}

func (g *ImageGrid) Z(c, r int) float64 {
	rr := g.rows - 1 - r
	return g.data[rr*g.cols+c]
}
