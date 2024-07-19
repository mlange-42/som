package table

import (
	"testing"

	"github.com/mlange-42/som/norm"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	tb := New([]string{"a", "b", "c"}, 5)

	assert.Equal(t, 15, len(tb.data))

	assert.Equal(t, 0, tb.Column("a"))
	assert.Equal(t, 1, tb.Column("b"))
	assert.Equal(t, 2, tb.Column("c"))

	for i := range tb.data {
		tb.data[i] = float64(i)
	}

	assert.Equal(t, 0, tb.index(0, 0))
	assert.Equal(t, 3, tb.index(1, 0))
	assert.Equal(t, 4, tb.index(1, 1))

	assert.Equal(t, 0.0, tb.Get(0, 0))
	assert.Equal(t, 3.0, tb.Get(1, 0))
	assert.Equal(t, 4.0, tb.Get(1, 1))
	assert.Equal(t, 7.0, tb.Get(2, 1))

	assert.Equal(t, []float64{0.0, 1.0, 2.0}, tb.GetRow(0))
	assert.Equal(t, []float64{12.0, 13.0, 14.0}, tb.GetRow(4))

	tb.Set(2, 3, 100.0)
	assert.Equal(t, 100.0, tb.Get(2, 3))
}

func TestNewTableFromData(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		columns := []string{"x", "y", "z"}
		data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
		table, err := NewWithData(columns, data)

		assert.NoError(t, err)
		assert.Equal(t, columns, table.columns)
		assert.Equal(t, 2, table.rows)
		assert.Equal(t, data, table.data)
	})

	t.Run("Empty columns", func(t *testing.T) {
		columns := []string{}
		data := []float64{1.0, 2.0, 3.0}
		_, err := NewWithData(columns, data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "columns length must be greater than zero")
	})

	t.Run("Data length not multiple of columns", func(t *testing.T) {
		columns := []string{"a", "b"}
		data := []float64{1.0, 2.0, 3.0}
		_, err := NewWithData(columns, data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data length 3 is not a multiple of columns length 2")
	})

	t.Run("Single column", func(t *testing.T) {
		columns := []string{"a"}
		data := []float64{1.0, 2.0, 3.0}
		table, err := NewWithData(columns, data)

		assert.NoError(t, err)
		assert.Equal(t, columns, table.columns)
		assert.Equal(t, 3, table.rows)
		assert.Equal(t, data, table.data)
	})

	t.Run("Empty data", func(t *testing.T) {
		columns := []string{"a", "b"}
		data := []float64{}
		table, err := NewWithData(columns, data)

		assert.NoError(t, err)
		assert.Equal(t, columns, table.columns)
		assert.Equal(t, 0, table.rows)
		assert.Equal(t, data, table.data)
	})
}
func TestNormalizeColumn(t *testing.T) {
	tb := New([]string{"a"}, 9)
	for i := range tb.data {
		tb.data[i] = float64(i)
	}

	minMaxNormalizer := &norm.Uniform{}
	minMaxNormalizer.SetArgs(0, 8)

	tb.NormalizeColumn(0, minMaxNormalizer)

	expected := []float64{0, 1 / 8.0, 2 / 8.0, 3 / 8.0, 4 / 8.0, 5 / 8.0, 6 / 8.0, 7 / 8.0, 8 / 8.0}
	assert.InDeltaSlice(t, expected, tb.data, 1e-6)

}

func BenchmarkTableGet(b *testing.B) {
	b.StopTimer()

	tb := New([]string{"a", "b", "c"}, 9)
	var v float64

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v = tb.Get(3, 1)
	}
	b.StopTimer()

	if v != 0.0 {
		b.Fatal("unexpected value")
	}
}

func BenchmarkTableGetRow(b *testing.B) {
	b.StopTimer()

	tb := New([]string{"a", "b", "c"}, 9)
	var v []float64

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v = tb.GetRow(3)
	}
	b.StopTimer()

	if len(v) != 3 {
		b.Fatal("unexpected value")
	}
}
