package som

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayer(t *testing.T) {
	l := NewLayer([]string{"a", "b", "c"}, Size{3, 2})

	assert.Equal(t, 18, len(l.data))

	assert.Equal(t, 0, l.Column("a"))
	assert.Equal(t, 1, l.Column("b"))
	assert.Equal(t, 2, l.Column("c"))

	for i := range l.data {
		l.data[i] = float64(i)
	}

	assert.Equal(t, 0.0, l.Get(0, 0, 0))
	assert.Equal(t, 1.0, l.Get(0, 0, 1))
	assert.Equal(t, 2.0, l.Get(0, 0, 2))
	assert.Equal(t, 3.0, l.Get(0, 1, 0))
	assert.Equal(t, 4.0, l.Get(0, 1, 1))
	assert.Equal(t, 5.0, l.Get(0, 1, 2))
	assert.Equal(t, 6.0, l.Get(1, 0, 0))
	assert.Equal(t, 7.0, l.Get(1, 0, 1))
	assert.Equal(t, 8.0, l.Get(1, 0, 2))
	assert.Equal(t, 9.0, l.Get(1, 1, 0))
	assert.Equal(t, 10.0, l.Get(1, 1, 1))
	assert.Equal(t, 11.0, l.Get(1, 1, 2))
	assert.Equal(t, 12.0, l.Get(2, 0, 0))
	assert.Equal(t, 13.0, l.Get(2, 0, 1))
	assert.Equal(t, 14.0, l.Get(2, 0, 2))
	assert.Equal(t, 15.0, l.Get(2, 1, 0))
	assert.Equal(t, 16.0, l.Get(2, 1, 1))
	assert.Equal(t, 17.0, l.Get(2, 1, 2))

	assert.Equal(t, []float64{0.0, 1.0, 2.0}, l.GetNode(0, 0))
	assert.Equal(t, []float64{3.0, 4.0, 5.0}, l.GetNode(0, 1))
	assert.Equal(t, []float64{15.0, 16.0, 17.0}, l.GetNode(2, 1))
}

func TestCoordsAt(t *testing.T) {
	l := NewLayer([]string{"a", "b", "c"}, Size{4, 3})

	testCases := []struct {
		name      string
		idx       int
		expectedX int
		expectedY int
	}{
		{"First element", 0, 0, 0},
		{"Last element", 11, 3, 2},
		{"Middle element", 5, 1, 1},
		{"Edge case - width", 3, 3, 0},
		{"Edge case - height", 8, 0, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, y := l.CoordsAt(tc.idx)
			assert.Equal(t, tc.expectedX, x, "Incorrect X coordinate")
			assert.Equal(t, tc.expectedY, y, "Incorrect Y coordinate")
		})
	}
}
