package conv

import (
	"math"
	"testing"

	"github.com/mlange-42/som/table"
	"github.com/stretchr/testify/assert"
)

func equalWithNaN(s1, s2 []float64) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if math.IsNaN(s1[i]) != math.IsNaN(s2[i]) {
			return false
		}
		if !math.IsNaN(s1[i]) && s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestClassesToTable(t *testing.T) {
	t.Run("String classes", func(t *testing.T) {
		classes := []string{"A", "B", "A", "C", "-", "A"}
		table, err := ClassesToTable(classes, nil, "-")
		assert.NoError(t, err)

		assert.Equal(t, 3, table.Columns())
		assert.Equal(t, 6, table.Rows())
		assert.Equal(t, []string{"A", "B", "C"}, table.ColumnNames())

		expectedData := []float64{
			1, 0, 0,
			0, 1, 0,
			1, 0, 0,
			0, 0, 1,
			math.NaN(), math.NaN(), math.NaN(),
			1, 0, 0,
		}
		assert.True(t, equalWithNaN(expectedData, table.Data()))
	})

	t.Run("String classes with columns", func(t *testing.T) {
		classes := []string{"A", "B", "A", "C", "B", "A"}
		table, err := ClassesToTable(classes, []string{"C", "A"}, "")
		assert.NoError(t, err)

		assert.Equal(t, 2, table.Columns())
		assert.Equal(t, 6, table.Rows())
		assert.Equal(t, []string{"C", "A"}, table.ColumnNames())

		expectedData := []float64{
			0, 1,
			math.NaN(), math.NaN(),
			0, 1,
			1, 0,
			math.NaN(), math.NaN(),
			0, 1,
		}
		assert.True(t, equalWithNaN(expectedData, table.Data()))
	})

	t.Run("Integer classes", func(t *testing.T) {
		classes := []int{1, 2, 1, 3, 2, 1, -1}
		table, err := ClassesToTable(classes, nil, -1)
		assert.NoError(t, err)

		assert.Equal(t, 3, table.Columns())
		assert.Equal(t, 7, table.Rows())
		assert.Equal(t, []string{"1", "2", "3"}, table.ColumnNames())

		expectedData := []float64{
			1, 0, 0,
			0, 1, 0,
			1, 0, 0,
			0, 0, 1,
			0, 1, 0,
			1, 0, 0,
			math.NaN(), math.NaN(), math.NaN(),
		}
		assert.True(t, equalWithNaN(expectedData, table.Data()))
	})

	t.Run("Single class", func(t *testing.T) {
		classes := []string{"A", "A", "A"}
		table, err := ClassesToTable(classes, nil, "")
		assert.NoError(t, err)

		assert.Equal(t, 1, table.Columns())
		assert.Equal(t, 3, table.Rows())
		assert.Equal(t, []string{"A"}, table.ColumnNames())

		expectedData := []float64{1, 1, 1}
		assert.Equal(t, expectedData, table.Data())
	})
}

func TestTableToClasses(t *testing.T) {
	t.Run("Basic case", func(t *testing.T) {
		table, err := table.NewWithData([]string{"A", "B", "C", "D"}, []float64{
			0.1, 0.2, 0.7, 0.0,
			0.8, 0.1, 0.0, 0.1,
			0.3, 0.3, 0.3, 0.1,
		})
		assert.Nil(t, err)

		columnNames, classes := TableToClasses(table)

		assert.Equal(t, []string{"A", "B", "C", "D"}, columnNames)
		assert.Equal(t, []int{2, 0, 0}, classes)
	})

	t.Run("Table with one row", func(t *testing.T) {
		table, err := table.NewWithData([]string{"X", "Y", "Z"}, []float64{
			0.1, 0.5, 0.4,
		})
		assert.Nil(t, err)

		columnNames, classes := TableToClasses(table)

		assert.Equal(t, []string{"X", "Y", "Z"}, columnNames)
		assert.Equal(t, []int{1}, classes)
	})

	t.Run("Table with equal values", func(t *testing.T) {
		table, err := table.NewWithData([]string{"P", "Q"}, []float64{
			0.5, 0.5,
			0.3, 0.3,
		})
		assert.Nil(t, err)

		columnNames, classes := TableToClasses(table)

		assert.Equal(t, []string{"P", "Q"}, columnNames)
		assert.Equal(t, []int{0, 0}, classes)
	})

	t.Run("Table with negative values", func(t *testing.T) {
		table, err := table.NewWithData([]string{"A", "B", "C"}, []float64{
			-0.1, -0.2, -0.3,
			-0.5, -0.1, -0.4,
		})
		assert.Nil(t, err)

		columnNames, classes := TableToClasses(table)

		assert.Equal(t, []string{"A", "B", "C"}, columnNames)
		assert.Equal(t, []int{0, 1}, classes)
	})
}
