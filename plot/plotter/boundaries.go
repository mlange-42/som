package plotter

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type GridBoundaries struct {
	plotter.GridXYZ
	draw.LineStyle
}

func NewGridBoundaries(grid plotter.GridXYZ) (*GridBoundaries, error) {
	style := plotter.DefaultLineStyle
	style.Width = 3
	return &GridBoundaries{
		GridXYZ:   grid,
		LineStyle: style,
	}, nil
}

func (h *GridBoundaries) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)

	dx := []int{0, 1, 0, -1}
	dy := []int{1, 0, -1, 0}

	cols, rows := h.GridXYZ.Dims()
	var pa vg.Path
	for i := 0; i < cols; i++ {
		left, right := h.getLeftRight(i, cols)
		for j := 0; j < rows; j++ {
			up, down := h.getUpDown(j, rows)
			vHere := h.GridXYZ.Z(i, j)
			for k := range dx {
				dxx := dx[k]
				dyy := dy[k]
				i2 := i + dxx
				j2 := j + dyy
				if i2 < 0 || i2 >= cols || j2 < 0 || j2 >= rows {
					continue
				}
				vThere := h.GridXYZ.Z(i2, j2)
				if vHere == vThere {
					continue
				}
				var x1, y1, x2, y2 font.Length
				if dxx < 0 {
					x1, y1 = trX(h.GridXYZ.X(i)+left), trY(h.GridXYZ.Y(j)+down)
					x2, y2 = trX(h.GridXYZ.X(i)+left), trY(h.GridXYZ.Y(j)+up)
				} else if dxx > 0 {
					x1, y1 = trX(h.GridXYZ.X(i)+right), trY(h.GridXYZ.Y(j)+down)
					x2, y2 = trX(h.GridXYZ.X(i)+right), trY(h.GridXYZ.Y(j)+up)
				} else if dyy < 0 {
					x1, y1 = trX(h.GridXYZ.X(i)+left), trY(h.GridXYZ.Y(j)+down)
					x2, y2 = trX(h.GridXYZ.X(i)+right), trY(h.GridXYZ.Y(j)+down)
				} else if dyy > 0 {
					x1, y1 = trX(h.GridXYZ.X(i)+left), trY(h.GridXYZ.Y(j)+up)
					x2, y2 = trX(h.GridXYZ.X(i)+right), trY(h.GridXYZ.Y(j)+up)
				}

				if !c.Contains(vg.Point{X: x1, Y: y1}) || !c.Contains(vg.Point{X: x2, Y: y2}) {
					continue
				}

				pa = pa[:0]
				pa.Move(vg.Point{X: x1, Y: y1})
				pa.Line(vg.Point{X: x2, Y: y2})

				c.SetLineStyle(h.LineStyle)
				c.Stroke(pa)
			}
		}
	}
}

func (h *GridBoundaries) getLeftRight(col, cols int) (left, right float64) {
	switch col {
	case 0:
		if cols == 1 {
			right = 0.5
		} else {
			right = (h.GridXYZ.X(1) - h.GridXYZ.X(0)) / 2
		}
		left = -right
	case cols - 1:
		right = (h.GridXYZ.X(cols-1) - h.GridXYZ.X(cols-2)) / 2
		left = -right
	default:
		right = (h.GridXYZ.X(col+1) - h.GridXYZ.X(col)) / 2
		left = -(h.GridXYZ.X(col) - h.GridXYZ.X(col-1)) / 2
	}
	return
}

func (h *GridBoundaries) getUpDown(row, rows int) (up, down float64) {
	switch row {
	case 0:
		if rows == 1 {
			up = 0.5
		} else {
			up = (h.GridXYZ.Y(1) - h.GridXYZ.Y(0)) / 2
		}
		down = -up
	case rows - 1:
		up = (h.GridXYZ.Y(rows-1) - h.GridXYZ.Y(rows-2)) / 2
		down = -up
	default:
		up = (h.GridXYZ.Y(row+1) - h.GridXYZ.Y(row)) / 2
		down = -(h.GridXYZ.Y(row) - h.GridXYZ.Y(row-1)) / 2
	}
	return
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (h *GridBoundaries) DataRange() (xMin, xMax, yMin, yMax float64) {
	c, r := h.GridXYZ.Dims()
	switch c {
	case 1: // Make a unit length when there is no neighbor.
		xMax = h.GridXYZ.X(0) + 0.5
		xMin = h.GridXYZ.X(0) - 0.5
	default:
		xMax = h.GridXYZ.X(c-1) + (h.GridXYZ.X(c-1)-h.GridXYZ.X(c-2))/2
		xMin = h.GridXYZ.X(0) - (h.GridXYZ.X(1)-h.GridXYZ.X(0))/2
	}
	switch r {
	case 1: // Make a unit length when there is no neighbor.
		yMax = h.GridXYZ.Y(0) + 0.5
		yMin = h.GridXYZ.Y(0) - 0.5
	default:
		yMax = h.GridXYZ.Y(r-1) + (h.GridXYZ.Y(r-1)-h.GridXYZ.Y(r-2))/2
		yMin = h.GridXYZ.Y(0) - (h.GridXYZ.Y(1)-h.GridXYZ.Y(0))/2
	}
	return xMin, xMax, yMin, yMax
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
func (h *GridBoundaries) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	c, r := h.GridXYZ.Dims()
	b := make([]plot.GlyphBox, 0, r*c)
	for i := 0; i < c; i++ {
		for j := 0; j < r; j++ {
			b = append(b, plot.GlyphBox{
				X: plt.X.Norm(h.GridXYZ.X(i)),
				Y: plt.Y.Norm(h.GridXYZ.Y(j)),
				Rectangle: vg.Rectangle{
					Min: vg.Point{X: -5, Y: -5},
					Max: vg.Point{X: +5, Y: +5},
				},
			})
		}
	}
	return b
}
