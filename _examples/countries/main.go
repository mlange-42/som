package main

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
)

func main() {
	path := "_examples/data/countries.csv"

	conf := som.SomConfig{
		Size:         layer.Size{10, 10},
		Neighborhood: &neighborhood.Gaussian{},
		Layers: []som.LayerDef{
			{
				Name:        "continent",
				Categorical: true,
				Metric:      &distance.Hamming{},
			},
			{
				Name: "scalars",
				Columns: []string{
					"child_mort_2010",
					"birth_p_1000",
					"log_GNI",
					"LifeExpectancy",
					"PopGrowth",
					"PopUrbanized",
					"PopGrowthUrb",
					"AdultLiteracy",
					"Income_low_40",
					"Income_high_20",
				},
				Norm: []norm.Normalizer{
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
					&norm.Gaussian{},
				},
				Metric: &distance.Euclidean{},
			},
		},
	}

	reader, err := csv.NewFileReader(path, ';', "-")
	if err != nil {
		panic(err)
	}
	tables, err := conf.PrepareTables(reader, true)
	if err != nil {
		panic(err)
	}

	s, err := som.New(&conf)
	if err != nil {
		panic(err)
	}

	trainingConf := som.TrainingConfig{
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Linear{Start: 10.0, End: 0.5},
	}

	trainer, err := som.NewTrainer(&s, tables, &trainingConf, rand.New(rand.NewSource(1)))
	if err != nil {
		panic(err)
	}

	trainer.Train(1000)

	for _, layer := range s.Layers() {
		fmt.Println("Layer", layer.Name())
		fmt.Println(layer.ToCSV(';'))

		if layer.IsCategorical() {
			classes, indices := conv.LayerToClasses(&layer)
			fmt.Println(classes)
			fmt.Println(indices)
		}
	}
}
