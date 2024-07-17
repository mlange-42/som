package yml

import (
	"bytes"
	"fmt"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"gopkg.in/yaml.v3"
)

type ymlLayer struct {
	Name        string
	Columns     []string `yaml:",flow,omitempty"`
	Metric      string
	Weight      float64
	Categorical bool      `yaml:",omitempty"`
	Data        []float64 `yaml:",flow,omitempty"`
}

type ymlConfig struct {
	Size         [2]int `yaml:",flow"`
	Layers       []ymlLayer
	Neighborhood string
}

func ToSomConfig(ymlData []byte) (*som.SomConfig, error) {
	reader := bytes.NewReader(ymlData)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	var yml ymlConfig
	err := decoder.Decode(&yml)
	if err != nil {
		return nil, err
	}

	neigh, ok := neighborhood.GetNeighborhood(yml.Neighborhood)
	if !ok {
		return nil, fmt.Errorf("unknown neighborhood: %s", yml.Neighborhood)
	}

	conf := som.SomConfig{
		Size:         layer.Size{Width: yml.Size[0], Height: yml.Size[1]},
		Layers:       []som.LayerDef{},
		Neighborhood: neigh,
	}
	for _, l := range yml.Layers {
		metric, ok := distance.GetMetric(l.Metric)
		if !ok {
			return nil, fmt.Errorf("unknown metric: %s", l.Metric)
		}
		if len(l.Data) > 0 && len(l.Data) != len(l.Columns)*yml.Size[0]*yml.Size[1] {
			return nil, fmt.Errorf("invalid data size for layer %s", l.Name)
		}
		conf.Layers = append(conf.Layers, som.LayerDef{
			Name:    l.Name,
			Columns: l.Columns,
			Metric:  metric,
			Weight:  l.Weight,
			Data:    l.Data,
		})
	}

	return &conf, nil
}

func ToYAML(som *som.Som) ([]byte, error) {
	yml := ymlConfig{
		Size:         [2]int{som.Size().Width, som.Size().Height},
		Layers:       []ymlLayer{},
		Neighborhood: som.Neighborhood().Name(),
	}
	for _, l := range som.Layers() {
		yml.Layers = append(yml.Layers, ymlLayer{
			Name:        l.Name(),
			Columns:     l.ColumnNames(),
			Metric:      l.Metric().Name(),
			Weight:      l.Weight(),
			Categorical: l.IsCategorical(),
			Data:        l.Data(),
		})
	}

	writer := bytes.Buffer{}
	encoder := yaml.NewEncoder(&writer)
	encoder.SetIndent(2)

	err := encoder.Encode(yml)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
