package som

import "fmt"

type Trainer struct {
	som   *Som
	table *Table
}

func NewTrainer(som *Som, table *Table) (*Trainer, error) {
	t := &Trainer{
		som:   som,
		table: table,
	}
	if !t.checkTable() {
		return nil, fmt.Errorf("table columns do not match SOM layer columns")
	}
	return t, nil
}

// checkTable checks that the table columns match the SOM layer columns.
// It returns true if the table columns match the SOM layer columns, false otherwise.
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
