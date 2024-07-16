package som

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSom(t *testing.T) {
	som := New(Size{2, 3}, []LayerDef{
		{Columns: []string{"x", "y"}},
		{Columns: []string{"a", "b", "c"}},
	})

	assert.Equal(t, 2, len(som.layers))
	assert.Equal(t, []int{0, 2}, som.offset)
	assert.Equal(t, Size{2, 3}, som.layers[0].size)
}
