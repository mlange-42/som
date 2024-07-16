package som

import "slices"

type Table struct {
	columns []string
	rows    int
	data    []float64
}

func NewTable(columns []string, rows int) Table {
	return Table{
		columns: columns,
		rows:    rows,
		data:    make([]float64, rows*len(columns)),
	}
}

func (t *Table) rowIndex(row int) int {
	return row * len(t.columns)
}

func (t *Table) index(row, col int) int {
	return row*len(t.columns) + col
}

func (t *Table) Column(col string) int {
	return slices.Index(t.columns, col)
}

func (t *Table) Get(row, col int) float64 {
	return t.data[t.index(row, col)]
}

func (t *Table) GetRow(row int) []float64 {
	idx := t.rowIndex(row)
	return t.data[idx : idx+len(t.columns)]
}
