package som

import (
	"fmt"
	"slices"
)

// Table represents a table of data with columns and rows.
type Table struct {
	columns []string  // The names of the columns in the table.
	rows    int       // The number of rows in the table.
	data    []float64 // The data values stored in row-major order.
}

// NewTable creates a new Table with the given column names and number of rows.
// The data slice is initialized to all zeros.
func NewTable(columns []string, rows int) Table {
	return Table{
		columns: columns,
		rows:    rows,
		data:    make([]float64, rows*len(columns)),
	}
}

// NewTableFromData creates a new Table from the given column names and data.
// If the length of the columns slice is zero, an error is returned.
// If the length of the data slice is not a multiple of the length of the columns slice, an error is returned.
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

// Column returns the index of the column with the given name. If the column is not found, -1 is returned.
func (t *Table) Column(col string) int {
	return slices.Index(t.columns, col)
}

// Get returns the value at the given row and column in the table.
func (t *Table) Get(row, col int) float64 {
	return t.data[t.index(row, col)]
}

// GetRow returns a slice containing the values for the given row in the table.
func (t *Table) GetRow(row int) []float64 {
	idx := t.rowIndex(row)
	return t.data[idx : idx+len(t.columns)]
}
