package som

type Trainer struct {
	som   *Som
	table *Table
}

func NewTrainer(som *Som, table *Table) *Trainer {
	t := &Trainer{
		som:   som,
		table: table,
	}
	t.checkTable()
	return t
}

func (t *Trainer) checkTable() bool {
	if len(t.table.columns) != t.som.offset[len(t.som.offset)-1]+len(t.som.layers[len(t.som.layers)-1].columns) {
		return false
	}
	for i := range t.som.layers {
		off := t.som.offset[i]
		for j, col := range t.som.layers[i].columns {
			if t.table.columns[j+off] != col {
				return false
			}
		}
	}
	return true
}
