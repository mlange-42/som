package som

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrainerCheckTable(t *testing.T) {
	som := New(Size{2, 3}, []LayerDef{
		{Columns: []string{"x", "y"}},
		{Columns: []string{"a", "b", "c"}},
	})

	t1 := NewTable([]string{"x", "y", "a", "b", "c"}, 5)
	tr := NewTrainer(&som, &t1)
	assert.True(t, tr.checkTable())

	t1 = NewTable([]string{"x", "y", "a", "b"}, 5)
	tr = NewTrainer(&som, &t1)
	assert.False(t, tr.checkTable())

	t1 = NewTable([]string{"x", "y", "a", "b", "c", "d"}, 5)
	tr = NewTrainer(&som, &t1)
	assert.False(t, tr.checkTable())

	t1 = NewTable([]string{"x", "a", "y", "b", "c"}, 5)
	tr = NewTrainer(&som, &t1)
	assert.False(t, tr.checkTable())
}
