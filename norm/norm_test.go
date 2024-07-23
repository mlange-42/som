package norm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTable struct {
	Mean float64
	Std  float64
	Min  float64
	Max  float64
}

func (m *mockTable) Range(col int) (min, max float64) {
	return m.Min, m.Max
}

func (m *mockTable) MeanStdDev(col int) (mean, std float64) {
	return m.Mean, m.Std
}

func TestGaussian(t *testing.T) {
	g := &Gaussian{mean: 10, std: 2}

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "gaussian", g.Name())
	})

	t.Run("Normalize", func(t *testing.T) {
		assert.InDelta(t, 0.5, g.Normalize(11), 0.001)
	})

	t.Run("DeNormalize", func(t *testing.T) {
		assert.InDelta(t, 11, g.DeNormalize(0.5), 0.001)
	})

	t.Run("Initialize", func(t *testing.T) {
		mock := mockTable{
			Mean: 2,
			Std:  1,
		}

		g.Initialize(&mock, 0)

		assert.Equal(t, 2.0, g.mean)
		assert.Equal(t, 1.0, g.std)
	})
}

func TestUniform(t *testing.T) {
	u := &Uniform{min: 0, max: 10}

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "uniform", u.Name())
	})

	t.Run("Normalize", func(t *testing.T) {
		assert.InDelta(t, 0.5, u.Normalize(5), 0.001)
	})

	t.Run("DeNormalize", func(t *testing.T) {
		assert.InDelta(t, 5, u.DeNormalize(0.5), 0.001)
	})

	t.Run("Initialize", func(t *testing.T) {
		mock := mockTable{
			Min: 2,
			Max: 6,
		}

		u.Initialize(&mock, 0)

		assert.Equal(t, 2.0, u.min)
		assert.Equal(t, 6.0, u.max)
	})
}

func TestNone(t *testing.T) {
	n := &Identity{}

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "none", n.Name())
	})

	t.Run("Normalize", func(t *testing.T) {
		assert.InDelta(t, 5, n.Normalize(5), 0.001)
	})

	t.Run("DeNormalize", func(t *testing.T) {
		assert.InDelta(t, 0.5, n.DeNormalize(0.5), 0.001)
	})

	t.Run("Initialize", func(t *testing.T) {
		mock := mockTable{}
		n.Initialize(&mock, 0)
		// No assertions needed for None as it doesn't modify any internal state
	})
}
