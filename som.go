package som

type LayerDef struct {
	Columns []string
}

type Som struct {
	size   Size
	layers []Layer
	offset []int
}

func New(size Size, layers []LayerDef) Som {
	lay := make([]Layer, len(layers))
	offset := make([]int, len(layers))
	off := 0
	for i, l := range layers {
		lay[i] = NewLayer(l.Columns, size)
		offset[i] = off
		off += len(l.Columns)
	}
	return Som{
		size:   size,
		layers: lay,
		offset: offset,
	}
}
