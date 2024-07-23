package plot

import (
	"image"
	"image/color"

	"github.com/mlange-42/som"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func Codes(s *som.Som, columns [][2]int, size image.Point) image.Image {
	legendHeight := 20

	codeHeight := (size.Y - legendHeight) / s.Size().Height
	codeWidth := size.X / s.Size().Width

	img := vgimg.NewWith(vgimg.UseWH(font.Length(size.X), font.Length(size.Y)), vgimg.UseDPI(72))
	dc := draw.New(img)

	style := draw.LineStyle{
		Color: color.Black,
		Width: 1,
	}

	w, h := s.Size().Width, s.Size().Height
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := draw.Crop(dc,
				font.Length(x*codeWidth), font.Length((x+1-w)*codeWidth),
				font.Length(y*codeHeight+legendHeight), font.Length((y+1-h)*codeHeight))
			c.StrokeLines(style, []vg.Point{
				{X: c.X(0.05), Y: c.Y(0.05)},
				{X: c.X(0.95), Y: c.Y(0.05)},
				{X: c.X(0.95), Y: c.Y(0.95)},
				{X: c.X(0.05), Y: c.Y(0.95)},
				{X: c.X(0.05), Y: c.Y(0.05)},
			})
		}
	}

	return img.Image()
}

func Code(s *som.Som, node int, columns [][2]int, size image.Point) {}
