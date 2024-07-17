package conv

import (
	"fmt"

	"github.com/mlange-42/som"
)

func ClassesToTable[T comparable](classes []T) ([]T, som.Table) {
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
	return classList, table
}
