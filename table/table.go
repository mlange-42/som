package table

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/mlange-42/som/norm"
)

type Reader interface {
	ReadColumns(columns []string) (*Table, error)
	ReadLabels(column string) ([]string, error)
}

// Table represents a table of data with columns and rows.
type Table struct {
	columns []string  // The names of the columns in the table.
	rows    int       // The number of rows in the table.
	data    []float64 // The data values stored in row-major order.
}

// New creates a new Table with the given column names and number of rows.
// The data slice is initialized to all zeros.
func New(columns []string, rows int) *Table {
	return &Table{
		columns: columns,
		rows:    rows,
		data:    make([]float64, rows*len(columns)),
	}
}

// NewWithData creates a new Table from the given column names and data.
// If the length of the columns slice is zero, an error is returned.
// If the length of the data slice is not a multiple of the length of the columns slice, an error is returned.
func NewWithData(columns []string, data []float64) (*Table, error) {
	if len(columns) == 0 {
		return nil, fmt.Errorf("columns length must be greater than zero")
	}
	if len(data)%len(columns) != 0 {
		return nil, fmt.Errorf("data length %d is not a multiple of columns length %d", len(data), len(columns))
	}
	return &Table{
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

func (t *Table) ColumnNames() []string {
	return t.columns
}

func (t *Table) Columns() int {
	return len(t.columns)
}

func (t *Table) Rows() int {
	return t.rows
}

// Get returns the value at the given row and column in the table.
func (t *Table) Get(row, col int) float64 {
	return t.data[t.index(row, col)]
}

// Set sets the value at the given row and column in the table.
func (t *Table) Set(row, col int, value float64) {
	t.data[t.index(row, col)] = value
}

// GetRow returns a slice containing the values for the given row in the table.
func (t *Table) GetRow(row int) []float64 {
	idx := t.rowIndex(row)
	return t.data[idx : idx+len(t.columns)]
}

func (t *Table) Data() []float64 {
	return t.data
}

func (t *Table) mean(col int) float64 {
	sum := 0.0
	count := 0
	for i := 0; i < t.Rows(); i++ {
		v := t.Get(i, col)
		if math.IsNaN(v) {
			continue
		}
		sum += v
		count++
	}
	return sum / float64(count)
}

func (t *Table) MeanStdDev(col int) (float64, float64) {
	mean := t.mean(col)
	sum := 0.0
	cnt := 0
	for i := 0; i < t.Rows(); i++ {
		v := t.Get(i, col)
		if math.IsNaN(v) {
			continue
		}
		diff := v - mean
		sum += diff * diff
		cnt++
	}
	return mean, math.Sqrt(sum / float64(cnt))
}

func (t *Table) Range(col int) (min, max float64) {
	min = math.Inf(1)
	max = math.Inf(-1)
	for i := 1; i < t.Rows(); i++ {
		v := t.Get(i, col)
		if math.IsNaN(v) {
			continue
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

func (t *Table) ToCSV(sep rune) string {
	b := strings.Builder{}
	cols := t.ColumnNames()
	for i, col := range cols {
		b.WriteString(col)
		if i < len(cols)-1 {
			b.WriteRune(sep)
		}
	}
	b.WriteRune('\n')
	for i := 0; i < t.Rows(); i++ {
		for j := 0; j < len(cols); j++ {
			b.WriteString(fmt.Sprintf("%f", t.Get(i, j)))
			if j < len(cols)-1 {
				b.WriteRune(sep)
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (t *Table) NormalizeColumn(col int, n norm.Normalizer) {
	for i := 0; i < t.Rows(); i++ {
		t.Set(i, col, n.Normalize(t.Get(i, col)))
	}
}
