package csv

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/mlange-42/som/table"
	"github.com/stretchr/testify/assert"
)

func TestReadColumns(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := "a,b,c\n1,2,3\n4,5,6"
		reader := strings.NewReader(input)
		columns := []string{"a", "c"}
		table, err := readColumns(reader, columns, ',', "")

		assert.NoError(t, err)
		assert.Equal(t, columns, table.ColumnNames())
		assert.Equal(t, 2, table.Rows())
		assert.Equal(t, []float64{1, 3, 4, 6}, table.Data())
	})

	t.Run("Missing column", func(t *testing.T) {
		input := "a,b,c\n1,2,3"
		reader := strings.NewReader(input)
		columns := []string{"a", "d"}
		_, err := readColumns(reader, columns, ',', "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "column \"d\" not found")
	})

	t.Run("No data value", func(t *testing.T) {
		input := "a,b,c\n1,NA,3\n4,5,NA"
		reader := strings.NewReader(input)
		columns := []string{"a", "b", "c"}
		table, err := readColumns(reader, columns, ',', "NA")

		assert.NoError(t, err)
		assert.Equal(t, 2, table.Rows())
		assert.True(t, math.IsNaN(table.Get(0, 1)))
		assert.True(t, math.IsNaN(table.Get(1, 2)))
	})

	t.Run("Invalid float", func(t *testing.T) {
		input := "a,b,c\n1,2,3\n4,invalid,6"
		reader := strings.NewReader(input)
		columns := []string{"a", "b", "c"}
		_, err := readColumns(reader, columns, ',', "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("Empty input", func(t *testing.T) {
		reader := strings.NewReader("")
		columns := []string{"a", "b"}
		_, err := readColumns(reader, columns, ',', "")

		assert.Error(t, err)
		assert.Equal(t, io.EOF, err)
	})

	t.Run("Custom delimiter", func(t *testing.T) {
		input := "a;b;c\n1;2;3\n4;5;6"
		reader := strings.NewReader(input)
		columns := []string{"b", "c"}
		table, err := readColumns(reader, columns, ';', "")

		assert.NoError(t, err)
		assert.Equal(t, columns, table.ColumnNames())
		assert.Equal(t, 2, table.Rows())
		assert.Equal(t, []float64{2, 3, 5, 6}, table.Data())
	})
}

func TestReadClasses(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := "a,b,c\nred,2,3\nblue,5,6\ngreen,8,9"
		reader := strings.NewReader(input)
		classes, err := readLabels(reader, "a", ',')

		assert.NoError(t, err)
		assert.Equal(t, []string{"red", "blue", "green"}, classes)
	})

	t.Run("Column not in first position", func(t *testing.T) {
		input := "x,y,z\n1,cat,3\n4,dog,6\n7,fish,9"
		reader := strings.NewReader(input)
		classes, err := readLabels(reader, "y", ',')

		assert.NoError(t, err)
		assert.Equal(t, []string{"cat", "dog", "fish"}, classes)
	})

	t.Run("Empty input", func(t *testing.T) {
		reader := strings.NewReader("")
		_, err := readLabels(reader, "a", ',')

		assert.Error(t, err)
		assert.Equal(t, io.EOF, err)
	})

	t.Run("Column not found", func(t *testing.T) {
		input := "a,b,c\n1,2,3\n4,5,6"
		reader := strings.NewReader(input)
		_, err := readLabels(reader, "d", ',')

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "column \"d\" not found")
	})

	t.Run("Custom delimiter", func(t *testing.T) {
		input := "a;b;c\napple;2;3\nbanana;5;6\ncherry;8;9"
		reader := strings.NewReader(input)
		classes, err := readLabels(reader, "a", ';')

		assert.NoError(t, err)
		assert.Equal(t, []string{"apple", "banana", "cherry"}, classes)
	})

	t.Run("Single column input", func(t *testing.T) {
		input := "class\nA\nB\nC"
		reader := strings.NewReader(input)
		classes, err := readLabels(reader, "class", ',')

		assert.NoError(t, err)
		assert.Equal(t, []string{"A", "B", "C"}, classes)
	})

	t.Run("Empty values in target column", func(t *testing.T) {
		input := "a,b,c\n1,,3\n,5,6\n7,,9"
		reader := strings.NewReader(input)
		classes, err := readLabels(reader, "b", ',')

		assert.NoError(t, err)
		assert.Equal(t, []string{"", "5", ""}, classes)
	})
}
func TestNewFileReader(t *testing.T) {
	t.Run("Valid file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("a,b,c\n1,2,3\n4,5,6")
		assert.NoError(t, err)
		tempFile.Close()

		reader, err := NewFileReader(tempFile.Name(), ',', "")
		assert.NoError(t, err)
		assert.NotNil(t, reader)
		assert.Equal(t, tempFile.Name(), reader.path)
		assert.Equal(t, "a,b,c\n1,2,3\n4,5,6", reader.text)
		assert.Equal(t, ',', reader.delim)
		assert.Equal(t, "", reader.noData)
	})

	t.Run("Non-existent file", func(t *testing.T) {
		_, err := NewFileReader("non_existent_file.csv", ',', "")
		assert.Error(t, err)
	})

	t.Run("Custom delimiter and noData", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("a;b;c\n1;NA;3\n4;5;NA")
		assert.NoError(t, err)
		tempFile.Close()

		reader, err := NewFileReader(tempFile.Name(), ';', "NA")
		assert.NoError(t, err)
		assert.NotNil(t, reader)
		assert.Equal(t, ';', reader.delim)
		assert.Equal(t, "NA", reader.noData)
	})
}

func TestFileReader_ReadColumns(t *testing.T) {
	t.Run("Read specific columns", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("a,b,c,d\n1,2,3,4\n5,6,7,8")
		assert.NoError(t, err)
		tempFile.Close()

		reader, err := NewFileReader(tempFile.Name(), ',', "")
		assert.NoError(t, err)

		table, err := reader.ReadColumns([]string{"a", "c"})
		assert.NoError(t, err)
		assert.NotNil(t, table)
		assert.Equal(t, []string{"a", "c"}, table.ColumnNames())
		assert.Equal(t, 2, table.Rows())
		assert.Equal(t, []float64{1, 3, 5, 7}, table.Data())
	})

	t.Run("Read with noData values", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("a,b,c\n1,N/A,3\n4,5,N/A")
		assert.NoError(t, err)
		tempFile.Close()

		reader, err := NewFileReader(tempFile.Name(), ',', "N/A")
		assert.NoError(t, err)

		table, err := reader.ReadColumns([]string{"a", "b", "c"})
		assert.NoError(t, err)
		assert.NotNil(t, table)
		assert.True(t, math.IsNaN(table.Get(0, 1)))
		assert.True(t, math.IsNaN(table.Get(1, 2)))
	})
}

func TestFileReader_ReadClasses(t *testing.T) {
	t.Run("Read classes from specific column", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("id,category,value\n1,A,10\n2,B,20\n3,A,30")
		assert.NoError(t, err)
		tempFile.Close()

		reader, err := NewFileReader(tempFile.Name(), ',', "")
		assert.NoError(t, err)

		classes, err := reader.ReadLabels("category")
		assert.NoError(t, err)
		assert.Equal(t, []string{"A", "B", "A"}, classes)
	})

	t.Run("Read classes with custom delimiter", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("id|type|value\n1|X|10\n2|Y|20\n3|Z|30")
		assert.NoError(t, err)
		tempFile.Close()

		reader, err := NewFileReader(tempFile.Name(), '|', "")
		assert.NoError(t, err)

		classes, err := reader.ReadLabels("type")
		assert.NoError(t, err)
		assert.Equal(t, []string{"X", "Y", "Z"}, classes)
	})

	t.Run("Read classes from non-existent column", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.csv")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("a,b,c\n1,2,3\n4,5,6")
		assert.NoError(t, err)
		tempFile.Close()

		reader, err := NewFileReader(tempFile.Name(), ',', "")
		assert.NoError(t, err)

		_, err = reader.ReadLabels("non_existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "column \"non_existent\" not found")
	})
}

func TestTableToCSV(t *testing.T) {
	t.Run("Empty table", func(t *testing.T) {
		tb := table.New([]string{"a", "b", "c"}, 0)
		var buf bytes.Buffer
		err := TableToCSV(tb, &buf, ',', "-")
		assert.NoError(t, err)
		expected := "a,b,c\n"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("Table with data", func(t *testing.T) {
		tb := table.New([]string{"x", "y"}, 2)
		tb.Set(0, 0, 1.5)
		tb.Set(0, 1, 2.0)
		tb.Set(1, 0, 3.5)
		tb.Set(1, 1, 4.0)
		var buf bytes.Buffer
		err := TableToCSV(tb, &buf, ',', "-")
		assert.NoError(t, err)
		expected := "x,y\n1.5,2\n3.5,4"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("Custom separator", func(t *testing.T) {
		tb := table.New([]string{"a", "b"}, 1)
		tb.Set(0, 0, 1.0)
		tb.Set(0, 1, 2.0)
		var buf bytes.Buffer
		err := TableToCSV(tb, &buf, ';', "-")
		assert.NoError(t, err)
		expected := "a;b\n1;2"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("Write error", func(t *testing.T) {
		tb := table.New([]string{"a"}, 1)
		tb.Set(0, 0, 1.0)
		mockWriter := &mockErrorWriter{err: fmt.Errorf("write error")}
		err := TableToCSV(tb, mockWriter, ',', "-")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "write error")
	})
}

type mockErrorWriter struct {
	err error
}

func (m *mockErrorWriter) Write(p []byte) (n int, err error) {
	return 0, m.err
}
