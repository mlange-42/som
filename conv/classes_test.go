package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassesToTable(t *testing.T) {
	t.Run("String classes", func(t *testing.T) {
		classes := []string{"A", "B", "A", "C", "B", "A"}
		classList, table := ClassesToTable(classes)

		assert.Equal(t, []string{"A", "B", "C"}, classList)
		assert.Equal(t, 3, table.Columns())
		assert.Equal(t, 6, table.Rows())
		assert.Equal(t, []string{"A", "B", "C"}, table.ColumnNames())

		expectedData := []float64{
			1, 0, 0,
			0, 1, 0,
			1, 0, 0,
			0, 0, 1,
			0, 1, 0,
			1, 0, 0,
		}
		assert.Equal(t, expectedData, table.Data())
	})

	t.Run("Integer classes", func(t *testing.T) {
		classes := []int{1, 2, 1, 3, 2, 1}
		classList, table := ClassesToTable(classes)

		assert.Equal(t, []int{1, 2, 3}, classList)
		assert.Equal(t, 3, table.Columns())
		assert.Equal(t, 6, table.Rows())
		assert.Equal(t, []string{"1", "2", "3"}, table.ColumnNames())

		expectedData := []float64{
			1, 0, 0,
			0, 1, 0,
			1, 0, 0,
			0, 0, 1,
			0, 1, 0,
			1, 0, 0,
		}
		assert.Equal(t, expectedData, table.Data())
	})

	t.Run("Single class", func(t *testing.T) {
		classes := []string{"A", "A", "A"}
		classList, table := ClassesToTable(classes)

		assert.Equal(t, []string{"A"}, classList)
		assert.Equal(t, 1, table.Columns())
		assert.Equal(t, 3, table.Rows())
		assert.Equal(t, []string{"A"}, table.ColumnNames())

		expectedData := []float64{1, 1, 1}
		assert.Equal(t, expectedData, table.Data())
	})
}
