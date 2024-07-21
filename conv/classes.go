package conv

import (
	"fmt"
	"math"

	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/table"
)

// ClassesToTable creates a table.Table from a slice of class labels.
// If the columns parameter is empty, the table will have a column for each unique class label.
// If the columns parameter is provided, the table will have a column for each provided class label,
// and an error will be returned if any duplicate class labels are found.
func ClassesToTable[T comparable](classes []T, columns []T) (*table.Table, error) {
	var classList []T
	var classNames []string
	classMap := map[T]int{}

	if len(columns) == 0 {
		for _, c := range classes {
			if _, ok := classMap[c]; !ok {
				classList = append(classList, c)
				classNames = append(classNames, fmt.Sprintf("%v", c))
				classMap[c] = len(classList) - 1
			}
		}
	} else {
		classList = columns
		classNames = make([]string, len(classList))
		for i, c := range classList {
			if _, ok := classMap[c]; !ok {
				classMap[c] = i
			} else {
				return nil, fmt.Errorf("duplicate class: %v", c)
			}
			classNames[i] = fmt.Sprintf("%v", c)
		}
	}

	data := make([]float64, len(classList)*len(classes))
	for i, c := range classes {
		idx, ok := classMap[c]
		if !ok {
			continue
		}
		data[i*len(classList)+idx] = 1
	}
	table, err := table.NewWithData(classNames, data)
	if err != nil {
		panic(err)
	}
	return table, nil
}

// ClassesToIndices converts a slice of class labels to a slice of class indices and a slice of unique class labels.
func ClassesToIndices[T comparable](classes []T) (columns []T, indices []int) {
	columns = []T{}
	classMap := map[T]int{}
	indices = make([]int, len(classes))

	for i, c := range classes {
		idx, ok := classMap[c]
		if !ok {
			idx = len(columns)
			classMap[c] = idx
			columns = append(columns, c)
		}
		indices[i] = idx
	}
	return columns, indices
}

// TableToClasses converts a table.Table into a slice of class labels and a slice of class indices.
// For each row in the table, the function finds the column with the maximum value and returns the index of that column as the class label.
// The function returns two slices: the first slice contains the column names (class labels), and the second slice contains the class indices for each row.
func TableToClasses(table *table.Table) ([]string, []int) {
	classes := make([]int, table.Rows())
	for i := 0; i < table.Rows(); i++ {
		row := table.GetRow(i)
		maxValue := math.Inf(-1)
		maxIndex := 0
		cols := table.Columns()
		for j := 0; j < cols; j++ {
			if row[j] > maxValue {
				maxValue = row[j]
				maxIndex = j
			}
		}
		classes[i] = maxIndex
	}
	return append([]string{}, table.ColumnNames()...), classes
}

// LayerToClasses converts a layer.Layer into a slice of class labels and a slice of class indices.
// For each node in the layer, the function finds the column with the maximum value and returns the index of that column as the class label.
// The function returns two slices: the first slice contains the column names (class labels), and the second slice contains the class indices for each node.
func LayerToClasses(l *layer.Layer) ([]string, []int) {
	units := l.Nodes()
	classes := make([]int, units)
	for i := 0; i < units; i++ {
		row := l.GetNodeAt(i)
		maxValue := math.Inf(-1)
		maxIndex := 0
		cols := l.ColumnNames()
		for j := 0; j < len(cols); j++ {
			if row[j] > maxValue {
				maxValue = row[j]
				maxIndex = j
			}
		}
		classes[i] = maxIndex
	}
	return append([]string{}, l.ColumnNames()...), classes
}
