package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"

	"github.com/mlange-42/som"
)

func ReadColumns(reader io.Reader, columns []string, delim rune, noData string) (*som.Table, error) {
	r := csv.NewReader(reader)
	r.ReuseRecord = true
	r.Comma = delim

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	indices := make([]int, len(columns))

	for i, c := range columns {
		idx := slices.Index(header, c)
		if idx == -1 {
			return nil, fmt.Errorf("column %q not found", c)
		}
		indices[i] = idx
	}

	data := []float64{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, idx := range indices {
			if record[idx] == noData {
				data = append(data, math.NaN())
				continue
			}
			v, err := strconv.ParseFloat(record[idx], 64)
			if err != nil {
				return nil, err
			}
			data = append(data, v)
		}
	}

	return som.NewTableFromData(columns, data)
}

func ReadClasses(reader io.Reader, column string, delim rune) ([]string, error) {
	r := csv.NewReader(reader)
	r.ReuseRecord = true
	r.Comma = delim

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	index := slices.Index(header, column)
	if index == -1 {
		return nil, fmt.Errorf("column %q not found", column)
	}

	data := []string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data = append(data, record[index])
	}

	return data, nil
}
