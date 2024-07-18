package plot

import (
	"image/color"
	"math/rand"

	"gonum.org/v1/plot/palette"
)

type RandomPalette struct {
	colors []color.Color
}

func NewRandomPalette(cols int) *RandomPalette {
	colors := make([]color.Color, cols)

	for i := range colors {
		col := palette.HSVA{
			H: float64(i) / float64(cols),
			S: rand.Float64()*0.5 + 0.5,
			V: rand.Float64()*0.2 + 0.8,
			A: 1,
		}
		colors[i] = color.NRGBAModel.Convert(col)
	}
	return &RandomPalette{colors: colors}
}

func (p *RandomPalette) Colors() []color.Color {
	return p.colors
}
