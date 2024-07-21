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

func TableToCSV(t *table.Table, writer io.Writer, sep rune, noData string) error {
	b := strings.Builder{}
	cols := t.ColumnNames()
	for i, col := range cols {
		b.WriteString(col)
		if i < len(cols)-1 {
			b.WriteRune(sep)
		}
	}
	b.WriteRune('\n')
	writer.Write([]byte(b.String()))
	b.Reset()

	for i := 0; i < t.Rows(); i++ {
		for j := 0; j < len(cols); j++ {
			v := t.Get(i, j)
			if math.IsNaN(v) {
				b.WriteString(noData)
			} else {
				b.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
			}
			if j < len(cols)-1 {
				b.WriteRune(sep)
			}
		}
		if i < t.Rows()-1 {
			b.WriteRune('\n')
		}
		_, err := writer.Write([]byte(b.String()))
		if err != nil {
			return err
		}
		b.Reset()
	}
	return nil
}

func TablesToCsv(tables []*table.Table, labelColumns []string, labels [][]string, writer io.Writer, delim rune, noData string) error {
	err := writeHeadersTables(writer, labelColumns, tables, delim)
	if err != nil {
		return err
	}

	builder := strings.Builder{}

	rows, err := checkAndCountTables(tables, labels)
	if err != nil {
		return err
	}

	for i := 0; i < rows; i++ {
		for j := range labels {
			builder.WriteString(labels[j][i])
			if i < len(labels)-1 || len(tables) > 0 {
				builder.WriteRune(delim)
			}
		}

		for j, tab := range tables {
			row := tab.GetRow(i)
			for k, v := range row {
				if math.IsNaN(v) {
					builder.WriteString(noData)
				} else {
					builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
				}
				if k < len(row)-1 || j < len(tables)-1 {
					builder.WriteRune(delim)
				}
			}
		}

		if i < rows-1 {
			builder.WriteRune('\n')
		}
		_, err := writer.Write([]byte(builder.String()))
		if err != nil {
			return err
		}
		builder.Reset()
	}

	return nil
}

func checkAndCountTables(tables []*table.Table, labels [][]string) (int, error) {
	rows := -1
	for _, t := range tables {
		if rows == -1 {
			rows = t.Rows()
		} else if rows != t.Rows() {
			return -1, fmt.Errorf("all tables and labels must have the same number of rows")
		}
	}
	for _, lab := range labels {
		if rows == -1 {
			rows = len(lab)
		} else if len(lab) != rows {
			return -1, fmt.Errorf("all tables and labels must have the same number of rows")
		}
	}

	return rows, nil
}

func writeHeadersTables(writer io.Writer, labelColumns []string, tables []*table.Table, delim rune) error {
	builder := strings.Builder{}

	for i, col := range labelColumns {
		builder.WriteString(col)
		if i < len(labelColumns)-1 || len(tables) > 0 {
			builder.WriteRune(delim)
		}
	}

	for i, tab := range tables {
		cols := tab.ColumnNames()
		for j, col := range cols {
			builder.WriteString(col)
			if j < len(cols)-1 || i < len(tables)-1 {
				builder.WriteRune(delim)
			}
		}
	}

	builder.WriteRune('\n')
	_, err := writer.Write([]byte(builder.String()))
	return err
}
