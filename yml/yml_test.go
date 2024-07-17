package yml

import (
	"testing"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/neighborhood"
	"github.com/stretchr/testify/assert"
)

func TestToSomConfig(t *testing.T) {
	t.Run("Valid YAML configuration", func(t *testing.T) {
		ymlData := []byte(`
size: [10, 8]
neighborhood: gaussian
layers:
  - name: layer1
    columns: [a, b, c]
    metric: euclidean
    weight: 1.0
  - name: layer2
    columns: [d, e]
    metric: manhattan
    weight: 0.5
`)

		config, err := ToSomConfig(ymlData)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, som.Size{Width: 10, Height: 8}, config.Size)
		assert.IsType(t, &neighborhood.Gaussian{}, config.Neighborhood)
		assert.Len(t, config.Layers, 2)
		assert.Equal(t, "layer1", config.Layers[0].Name)
		assert.Equal(t, []string{"a", "b", "c"}, config.Layers[0].Columns)
		assert.Equal(t, &distance.Euclidean{}, config.Layers[0].Metric)
		assert.Equal(t, 1.0, config.Layers[0].Weight)
		assert.Equal(t, "layer2", config.Layers[1].Name)
		assert.Equal(t, []string{"d", "e"}, config.Layers[1].Columns)
		assert.Equal(t, &distance.Manhattan{}, config.Layers[1].Metric)
		assert.Equal(t, 0.5, config.Layers[1].Weight)
	})

	t.Run("Invalid YAML syntax", func(t *testing.T) {
		ymlData := []byte(`
size: [10, 8
neighborhood: gaussian
`)

		config, err := ToSomConfig(ymlData)

		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("Unknown neighborhood", func(t *testing.T) {
		ymlData := []byte(`
size: [10, 8]
neighborhood: unknown
layers:
  - columns: [a, b, c]
    metric: euclidean
    weight: 1.0
`)

		config, err := ToSomConfig(ymlData)

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "unknown neighborhood: unknown")
	})

	t.Run("Unknown metric", func(t *testing.T) {
		ymlData := []byte(`
size: [10, 8]
neighborhood: gaussian
layers:
  - columns: [a, b, c]
    metric: unknown
    weight: 1.0
`)

		config, err := ToSomConfig(ymlData)

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "unknown metric: unknown")
	})

	t.Run("Empty configuration", func(t *testing.T) {
		ymlData := []byte(``)

		config, err := ToSomConfig(ymlData)

		assert.Error(t, err)
		assert.Nil(t, config)
	})
}
