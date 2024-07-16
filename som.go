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

// CheckTable checks if the provided Table matches the structure of the Som.
// It verifies that the total number of columns in the Table matches the
// total number of columns defined in the Som's layers. It also checks that
// the column names in the Table match the column names defined in the Som's
// layers.
func (s *Som) CheckTable(t *Table) bool {
	if len(t.columns) != s.offset[len(s.offset)-1]+len(s.layers[len(s.layers)-1].columns) {
		return false
	}
	for i := range s.layers {
		off := s.offset[i]
		for j, col := range s.layers[i].columns {
			if t.columns[j+off] != col {
				return false
			}
		}
	}
	return true
}
