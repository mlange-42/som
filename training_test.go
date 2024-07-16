package som

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrainer(t *testing.T) {
	som := New(Size{2, 3}, []LayerDef{
		{Columns: []string{"x", "y"}},
		{Columns: []string{"a", "b", "c"}},
	})

	t1 := NewTable([]string{"x", "y", "a", "b", "c"}, 5)
	_, err := NewTrainer(&som, &t1)
	assert.Nil(t, err)

	t1 = NewTable([]string{"x", "y", "a", "b"}, 5)
	_, err = NewTrainer(&som, &t1)
	assert.NotNil(t, err)

	t1 = NewTable([]string{"x", "y", "a", "b", "c", "d"}, 5)
	_, err = NewTrainer(&som, &t1)
	assert.NotNil(t, err)

	t1 = NewTable([]string{"x", "a", "y", "b", "c"}, 5)
	_, err = NewTrainer(&som, &t1)
	assert.NotNil(t, err)
}
