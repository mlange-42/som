package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/mlange-42/som/table"
)

type StringReader struct {
	text   string
	delim  rune
	noData string
}

func NewStringReader(text string, delim rune, noData string) *StringReader {
	return &StringReader{text, delim, noData}
}

func (s *StringReader) ReadColumns(columns []string) (*table.Table, error) {
	return readColumns(strings.NewReader(s.text), columns, s.delim, s.noData)
}

func (s *StringReader) ReadLabels(column string) ([]string, error) {
	return readLabels(strings.NewReader(s.text), column, s.delim)
}

type FileReader struct {
	path   string
	text   string
	delim  rune
	noData string
}

func NewFileReader(path string, delim rune, noData string) (*FileReader, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &FileReader{path, string(text), delim, noData}, nil
}

func (f *FileReader) ReadColumns(columns []string) (*table.Table, error) {
	return readColumns(strings.NewReader(f.text), columns, f.delim, f.noData)
}

func (f *FileReader) ReadLabels(column string) ([]string, error) {
	return readLabels(strings.NewReader(f.text), column, f.delim)
}

func readColumns(reader io.Reader, columns []string, delim rune, noData string) (*table.Table, error) {
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

	return table.NewWithData(columns, data)
}

func readLabels(reader io.Reader, column string, delim rune) ([]string, error) {
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
