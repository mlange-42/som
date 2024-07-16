package som

import (
	"fmt"
	"slices"
)

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

func NewTableFromData(columns []string, data []float64) (Table, error) {
	if len(columns) == 0 {
		return Table{}, fmt.Errorf("columns length must be greater than zero")
	}
	if len(data)%len(columns) != 0 {
		return Table{}, fmt.Errorf("data length %d is not a multiple of columns length %d", len(data), len(columns))
	}
	return Table{
		columns: columns,
		rows:    len(data) / len(columns),
		data:    data,
	}, nil
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
