package yml

import (
	"bytes"
	"fmt"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/neighborhood"
	"gopkg.in/yaml.v3"
)

type ymlLayer struct {
	Name    string
	Columns []string
	Metric  string
	Weight  float64
}

type ymlConfig struct {
	Size         [2]int
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
		Size:         som.Size{Width: yml.Size[0], Height: yml.Size[1]},
		Layers:       []som.LayerDef{},
		Neighborhood: neigh,
	}
	for _, l := range yml.Layers {
		metric, ok := distance.GetMetric(l.Metric)
		if !ok {
			return nil, fmt.Errorf("unknown metric: %s", l.Metric)
		}
		conf.Layers = append(conf.Layers, som.LayerDef{
			Name:    l.Name,
			Columns: l.Columns,
			Metric:  metric,
			Weight:  l.Weight,
		})
	}

	return &conf, nil
}
