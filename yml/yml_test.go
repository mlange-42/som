package yml

import (
	"testing"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
	"github.com/stretchr/testify/assert"
)

func TestToSomConfig(t *testing.T) {
	t.Run("Valid YAML configuration", func(t *testing.T) {
		ymlData := []byte(`
size: [4, 3]
neighborhood: gaussian
layers:
  - name: layer1
    columns: [a, b, c]
    norm: [gaussian 0 1, uniform, none]
    metric: euclidean
    weight: 1.0
  - name: layer2
    columns: [d, e]
    metric: manhattan
    weight: 0.5
    data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
`)

		config, err := ToSomConfig(ymlData)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, layer.Size{Width: 4, Height: 3}, config.Size)
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

		gauss := norm.Gaussian{}
		gauss.SetArgs(0, 1)

		assert.Equal(t, &gauss, config.Layers[0].Norm[0])
		assert.Equal(t, &norm.Uniform{}, config.Layers[0].Norm[1])
		assert.Equal(t, &norm.None{}, config.Layers[0].Norm[2])
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

	t.Run("Invalid data size", func(t *testing.T) {
		ymlData := []byte(`
size: [4, 3]
neighborhood: gaussian
layers:
  - name: Layer A
    columns: [a, b, c]
    weight: 1.0
    metric: euclidean
    data: [0, 0, 0, 0]
`)

		config, err := ToSomConfig(ymlData)

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "invalid data size for layer Layer A")
	})

	t.Run("Empty configuration", func(t *testing.T) {
		ymlData := []byte(``)

		config, err := ToSomConfig(ymlData)

		assert.Error(t, err)
		assert.Nil(t, config)
	})
}
func TestToYAML(t *testing.T) {
	ymlData := []byte(`
size: [4, 3]
neighborhood: gaussian
layers:
- name: layer1
  columns: [a, b, c]
  metric: euclidean
  weight: 1.0
- name: layer2
  columns: [d, e]
  norm: [gaussian 0 1, uniform -0.01 0.01]
  metric: manhattan
  weight: 0.5
`)

	config, err := ToSomConfig(ymlData)
	assert.NoError(t, err)

	s, err := som.New(config)
	assert.NoError(t, err)

	result, err := ToYAML(&s)
	assert.NoError(t, err)

	expected := `size: [4, 3]
neighborhood: gaussian
layers:
  - name: layer1
    columns: [a, b, c]
    metric: euclidean
    weight: 1
    data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
  - name: layer2
    columns: [d, e]
    norm: [gaussian 0 1, uniform -0.01 0.01]
    metric: manhattan
    weight: 0.5
    data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
`

	assert.Equal(t, expected, string(result))
}
