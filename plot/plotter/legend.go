package plotter

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Legend struct {
	TextStyle      text.Style
	Padding        vg.Length
	Top, Left      bool
	YPosition      float64
	XOffs, YOffs   vg.Length
	ThumbnailWidth vg.Length
	Columns        int
	entries        []legendEntry
}

type legendEntry struct {
	text   string
	thumbs []plot.Thumbnailer
}

func NewLegend() Legend {
	return newLegend(plot.DefaultTextHandler)
}

func newLegend(hdlr text.Handler) Legend {
	return Legend{
		Columns:        1,
		YPosition:      draw.PosBottom,
		ThumbnailWidth: vg.Points(20),
		TextStyle: text.Style{
			Font:    font.From(plot.DefaultFont, 12),
			Handler: hdlr,
		},
	}
}

func (l *Legend) Draw(c draw.Canvas) {
	sty := l.TextStyle
	descent := sty.FontExtents().Descent
	em := sty.Rectangle(" ")

	entryHeight, textWidth := l.entryHeight(), l.textWidth()
	rows := int(math.Ceil(float64(len(l.entries)) / float64(l.Columns)))

	iconOffset := font.Length(0)
	textOffset := l.ThumbnailWidth + em.Max.X
	xOrigin := c.Min.X
	yOrigin := c.Max.Y - font.Length(rows)*(entryHeight+l.Padding)
	if !l.Left {
		iconOffset = textWidth
		textOffset = textWidth - em.Max.X
		xOrigin = c.Max.X - font.Length(l.Columns)*(l.ThumbnailWidth+textWidth)
		sty.XAlign--
	}
	if !l.Top {
		yOrigin = c.Min.Y
	}
	xOrigin += l.XOffs
	yOrigin += l.YOffs

	icon := &draw.Canvas{
		Canvas: c.Canvas,
	}

	if l.YPosition < draw.PosBottom || draw.PosTop < l.YPosition {
		panic("plot: invalid vertical offset for the legend's entries")
	}
	yoff := vg.Length(l.YPosition-draw.PosBottom) / 2
	yoff += descent

	w := textWidth + l.ThumbnailWidth
	for i, e := range l.entries {
		row := i / l.Columns
		col := i % l.Columns
		x := vg.Length(col)*w + xOrigin
		y := vg.Length(rows-1-row)*(entryHeight+l.Padding) + yOrigin

		icon.Rectangle = vg.Rectangle{
			Min: vg.Point{X: x + iconOffset, Y: y},
			Max: vg.Point{X: x + iconOffset + l.ThumbnailWidth, Y: y + entryHeight},
		}

		for _, t := range e.thumbs {
			t.Thumbnail(icon)
		}
		yoffs := (entryHeight - descent - sty.Rectangle(e.text).Max.Y) / 2
		yoffs += yoff

		c.FillText(sty, vg.Point{X: x + textOffset, Y: icon.Min.Y + yoffs}, e.text)
	}
}

func (l *Legend) AdjustColumns(maxWidth font.Length) {
	w := l.ThumbnailWidth + l.textWidth() + font.Millimeter
	l.Columns = min(len(l.entries), int(maxWidth/w))
}

func (l *Legend) Rectangle(c draw.Canvas) vg.Rectangle {
	entryHeight := l.entryHeight()
	rows := int(math.Ceil(float64(len(l.entries)) / float64(l.Columns)))

	height := entryHeight*vg.Length(rows) + l.Padding*vg.Length(rows-1)
	width := (l.ThumbnailWidth + l.textWidth()) * vg.Length(l.Columns)

	var r vg.Rectangle
	if l.Left {
		r.Max.X = c.Min.X + width
		r.Min.X = c.Min.X
	} else {
		r.Max.X = c.Max.X
		r.Min.X = c.Max.X - width
	}
	if l.Top {
		r.Max.Y = c.Max.Y
		r.Min.Y = c.Max.Y - height
	} else {
		r.Max.Y = c.Min.Y + height
		r.Min.Y = c.Min.Y
	}
	return r
}

func (l *Legend) entryHeight() (height vg.Length) {
	for _, e := range l.entries {
		if h := l.TextStyle.Rectangle(e.text).Max.Y; h > height {
			height = h
		}
	}
	return
}

func (l *Legend) textWidth() (width vg.Length) {
	for _, e := range l.entries {
		t := e.text + " "
		if l.Columns > 1 {
			t += " "
		}
		if w := l.TextStyle.Rectangle(t).Max.X; w > width {
			width = w
		}
	}
	return
}

func (l *Legend) Add(name string, thumbs ...plot.Thumbnailer) {
	l.entries = append(l.entries, legendEntry{text: name, thumbs: thumbs})
}
