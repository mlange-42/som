package som

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	tb := NewTable([]string{"a", "b", "c"}, 5)

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
}
