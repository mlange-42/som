package norm

import (
	"testing"

	"github.com/mlange-42/som/table"
	"github.com/stretchr/testify/assert"
)

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
		mockTable, err := table.NewWithData([]string{"x"}, []float64{
			2.0, 2.0, 2.0,
		})
		assert.NoError(t, err)

		g.Initialize(mockTable, 0)

		assert.Equal(t, 2.0, g.mean)
		assert.Equal(t, 0.0, g.std)
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
		mockTable, err := table.NewWithData([]string{"x"}, []float64{
			2.0, 4.0, 6.0,
		})
		assert.NoError(t, err)

		u.Initialize(mockTable, 0)

		assert.Equal(t, 2.0, u.min)
		assert.Equal(t, 6.0, u.max)
	})
}
