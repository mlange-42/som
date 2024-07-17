package conv

import (
	"fmt"
	"math"

	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/table"
)

func ClassesToTable[T comparable](classes []T) *table.Table {
	classList := []T{}
	classNames := []string{}
	classMap := map[T]int{}

	for _, c := range classes {
		if _, ok := classMap[c]; !ok {
			classList = append(classList, c)
			classNames = append(classNames, fmt.Sprintf("%v", c))
			classMap[c] = len(classList) - 1
		}
	}

	data := make([]float64, len(classList)*len(classes))
	for i, c := range classes {
		data[i*len(classList)+classMap[c]] = 1
	}
	table, err := table.NewFromData(classNames, data)
	if err != nil {
		panic(err)
	}
	return table
}

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
