package plot

import (
	"image/color"
)

type RandomPalette struct {
	colors []color.Color
}

func NewRandomPalette(cols []color.Color, count int) *RandomPalette {
	colors := make([]color.Color, count)
	for i := range colors {
		colors[i] = cols[i%len(cols)]
	}
	/*for i := range colors {
		col := palette.HSVA{
			H: float64(i) / float64(cols),
			S: rand.Float64()*0.5 + 0.5,
			V: rand.Float64()*0.2 + 0.8,
			A: 1,
		}
		colors[i] = color.NRGBAModel.Convert(col)
	}*/
	return &RandomPalette{colors: colors}
}

func (p *RandomPalette) Colors() []color.Color {
	return p.colors
}

type WhitePalette struct{}

func (p *WhitePalette) Colors() []color.Color {
	return []color.Color{color.White}
}
