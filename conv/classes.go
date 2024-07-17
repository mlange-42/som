package conv

import (
	"fmt"
	"math"

	"github.com/mlange-42/som"
)

func ClassesToTable[T comparable](classes []T) *som.Table {
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
	table, err := som.NewTableFromData(classNames, data)
	if err != nil {
		panic(err)
	}
	return table
}

func TableToClasses(table *som.Table) ([]string, []int) {
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
