package plot

import (
	"github.com/mlange-42/som"
	"github.com/mlange-42/som/layer"
)

type SomLayerGrid struct {
	Som    *som.Som
	Layer  int
	Column int
}

func (g *SomLayerGrid) Dims() (c, r int) {
	return g.Som.Size().Width, g.Som.Size().Height
}

func (g *SomLayerGrid) Z(c, r int) float64 {
	l := g.Som.Layers()[g.Layer]
	v := l.Get(c, r, g.Column)
	return l.Normalizers()[g.Column].DeNormalize(v)
}

func (g *SomLayerGrid) X(c int) float64 {
	return float64(c)
}

func (g *SomLayerGrid) Y(r int) float64 {
	return float64(r)
}

type ClassesGrid struct {
	Size    layer.Size
	Indices []int
}

func (g *ClassesGrid) Dims() (c, r int) {
	return g.Size.Width, g.Size.Height
}

func (g *ClassesGrid) Z(c, r int) float64 {
	idx := r + c*g.Size.Height
	return float64(g.Indices[idx])
}

func (g *ClassesGrid) X(c int) float64 {
	return float64(c)
}

func (g *ClassesGrid) Y(r int) float64 {
	return float64(r)
}

type UMatrixGrid struct {
	UMatrix [][]float64
}

func (g *UMatrixGrid) Dims() (c, r int) {
	return len(g.UMatrix[0]), len(g.UMatrix)
}

func (g *UMatrixGrid) Z(c, r int) float64 {
	return g.UMatrix[r][c]
}

func (g *UMatrixGrid) X(c int) float64 {
	return float64(c) / 2
}

func (g *UMatrixGrid) Y(r int) float64 {
	return float64(r) / 2
}
