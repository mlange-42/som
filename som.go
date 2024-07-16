package som

type LayerDef struct {
	Columns []string
	Weight  float64
}

type Som struct {
	size   Size
	layers []Layer
	offset []int
	weight []float64
}

func New(size Size, layers []LayerDef) Som {
	lay := make([]Layer, len(layers))
	weight := make([]float64, len(layers))
	offset := make([]int, len(layers))
	off := 0
	for i, l := range layers {
		lay[i] = NewLayer(l.Columns, size)
		offset[i] = off
		off += len(l.Columns)

		weight[i] = l.Weight
		if weight[i] == 0 {
			weight[i] = 1
		}
	}
	return Som{
		size:   size,
		layers: lay,
		offset: offset,
		weight: weight,
	}
}
