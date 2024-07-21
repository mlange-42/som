package layer

import (
	"testing"

	"github.com/mlange-42/som/distance"
	"github.com/stretchr/testify/assert"
)

func TestLayer(t *testing.T) {
	l, err := New("L1", []string{"a", "b", "c"}, nil, Size{3, 2}, &distance.Manhattan{}, 1.0, false)
	assert.NoError(t, err)

	assert.Equal(t, 18, len(l.weights))

	assert.Equal(t, 0, l.Column("a"))
	assert.Equal(t, 1, l.Column("b"))
	assert.Equal(t, 2, l.Column("c"))

	for i := range l.weights {
		l.weights[i] = float64(i)
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

func BenchmarkLayerGet(b *testing.B) {
	b.StopTimer()

	l, err := New("L1", []string{"a", "b", "c"}, nil, Size{3, 2}, &distance.Manhattan{}, 1.0, false)
	if err != nil {
		b.Fatal(err)
	}
	var v float64

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v = l.Get(0, 1, 2)
	}
	b.StopTimer()

	if v != 0.0 {
		b.Fatal("unexpected value")
	}
}

func BenchmarkLayerGetNode(b *testing.B) {
	b.StopTimer()

	l, err := New("L1", []string{"a", "b", "c"}, nil, Size{3, 2}, &distance.Manhattan{}, 1.0, false)
	if err != nil {
		b.Fatal(err)
	}
	var v []float64

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v = l.GetNode(0, 1)
	}
	b.StopTimer()

	if len(v) != 3 {
		b.Fatal("unexpected value")
	}
}

func BenchmarkLayerCoordsAt(b *testing.B) {
	b.StopTimer()

	l, err := New("L1", []string{"a", "b", "c"}, nil, Size{3, 2}, &distance.Manhattan{}, 1.0, false)
	if err != nil {
		b.Fatal(err)
	}

	var x, y int

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		x, y = l.CoordsAt(3)
	}

	b.StopTimer()
	if x != 1 || y != 1 {
		b.Fatal("unexpected value")
	}
}
