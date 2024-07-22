package yml

import (
	"bytes"
	"fmt"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
	"gopkg.in/yaml.v3"
)

type ymlLayer struct {
	Name        string
	Columns     []string `yaml:",flow,omitempty"`
	Norm        []string `yaml:",flow,omitempty"`
	Metric      string
	Weight      float64   `yaml:",omitempty"`
	Categorical bool      `yaml:",omitempty"`
	Data        []float64 `yaml:",flow,omitempty"`
}

type ymlSom struct {
	Size         [2]int `yaml:",flow"`
	Neighborhood string
	Metric       string
	Layers       []*ymlLayer
}

type ymlTraining struct {
	Epochs      int
	Alpha       string `yaml:",omitempty"`
	Radius      string `yaml:",omitempty"`
	WeightDecay string `yaml:"weight-decay,omitempty"`
	Lambda      float64
}

type ymlConfig struct {
	Som      ymlSom
	Training *ymlTraining `yaml:",omitempty"`
}

func ToSomConfig(ymlData []byte) (*som.SomConfig, *som.TrainingConfig, error) {
	reader := bytes.NewReader(ymlData)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	var yml ymlConfig
	err := decoder.Decode(&yml)
	if err != nil {
		return nil, nil, err
	}

	neigh, ok := neighborhood.GetNeighborhood(yml.Som.Neighborhood)
	if !ok {
		return nil, nil, fmt.Errorf("unknown neighborhood: %s", yml.Som.Neighborhood)
	}
	metric, ok := neighborhood.GetMetric(yml.Som.Metric)
	if !ok {
		return nil, nil, fmt.Errorf("unknown neighborhood metric: %s", yml.Som.Metric)
	}

	conf := som.SomConfig{
		Size:         layer.Size{Width: yml.Som.Size[0], Height: yml.Som.Size[1]},
		Layers:       []*som.LayerDef{},
		Neighborhood: neigh,
		MapMetric:    metric,
	}
	for _, l := range yml.Som.Layers {
		lay, err := createLayer(&yml.Som, l)
		if err != nil {
			return nil, nil, err
		}
		conf.Layers = append(conf.Layers, lay)
	}

	var training *som.TrainingConfig

	if yml.Training != nil {
		alpha, err := decay.FromString(yml.Training.Alpha)
		if err != nil {
			return nil, nil, err
		}
		radius, err := decay.FromString(yml.Training.Radius)
		if err != nil {
			return nil, nil, err
		}

		var wtDecay decay.Decay
		if yml.Training.WeightDecay != "" {
			wtDecay, err = decay.FromString(yml.Training.WeightDecay)
			if err != nil {
				return nil, nil, err
			}
		}

		training = &som.TrainingConfig{
			Epochs:             yml.Training.Epochs,
			LearningRate:       alpha,
			NeighborhoodRadius: radius,
			WeightDecay:        wtDecay,
			ViSomLambda:        yml.Training.Lambda,
		}
	}

	return &conf, training, nil
}

func createLayer(s *ymlSom, l *ymlLayer) (*som.LayerDef, error) {
	metric, ok := distance.GetMetric(l.Metric)
	if !ok {
		return nil, fmt.Errorf("unknown metric: %s", l.Metric)
	}
	if len(l.Data) > 0 && len(l.Data) != len(l.Columns)*s.Size[0]*s.Size[1] {
		return nil, fmt.Errorf("invalid data size for layer %s", l.Name)
	}

	if len(l.Norm) > 1 && len(l.Norm) != len(l.Columns) {
		return nil, fmt.Errorf("invalid number of normalizers for layer %s; must be zero, one or number of columns", l.Name)
	}

	norms := make([]norm.Normalizer, len(l.Columns))
	for i := range norms {
		var err error
		if i >= len(l.Norm) {
			if len(l.Norm) == 0 {
				norms[i] = &norm.None{}
				continue
			}
			norms[i], err = norm.FromString(l.Norm[0])
			if err != nil {
				return nil, err
			}
			continue
		}
		norms[i], err = norm.FromString(l.Norm[i])
		if err != nil {
			return nil, err
		}
	}

	return &som.LayerDef{
		Name:        l.Name,
		Columns:     l.Columns,
		Norm:        norms,
		Metric:      metric,
		Weight:      l.Weight,
		Weights:     l.Data,
		Categorical: l.Categorical,
	}, nil
}

func ToYAML(som *som.Som) ([]byte, error) {
	yml := ymlSom{
		Size:         [2]int{som.Size().Width, som.Size().Height},
		Layers:       []*ymlLayer{},
		Neighborhood: som.Neighborhood().Name(),
		Metric:       som.MapMetric().Name(),
	}
	for _, l := range som.Layers() {
		norms := make([]string, len(l.Normalizers()))
		allNone := true
		for i, n := range l.Normalizers() {
			norms[i] = norm.ToString(n)
			if _, ok := n.(*norm.None); !ok {
				allNone = false
			}
		}
		if allNone {
			norms = nil
		}

		weight := l.Weight()
		if weight == 0.0 {
			weight = -1
		} else if weight == 1.0 {
			weight = 0
		}

		yml.Layers = append(yml.Layers, &ymlLayer{
			Name:        l.Name(),
			Columns:     l.ColumnNames(),
			Norm:        norms,
			Metric:      l.Metric().Name(),
			Weight:      weight,
			Categorical: l.IsCategorical(),
			Data:        l.Weights(),
		})
	}

	ymlConf := ymlConfig{
		Som:      yml,
		Training: nil,
	}

	writer := bytes.Buffer{}
	encoder := yaml.NewEncoder(&writer)
	encoder.SetIndent(2)

	err := encoder.Encode(ymlConf)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}
