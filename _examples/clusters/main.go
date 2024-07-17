package main

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/neighborhood"
)

func main() {
	rng := rand.New(rand.NewSource(1))

	somParams := som.SomConfig{
		Size:         som.Size{Width: 8, Height: 6},
		Layers:       []som.LayerDef{{Columns: []string{"x", "y"}}},
		Neighborhood: &neighborhood.Linear{},
	}
	s := som.New(&somParams)

	rows := 250
	data := generateData(rows, 2)
	table, err := som.NewTableFromData([]string{"x", "y"}, data)
	if err != nil {
		panic(err)
	}

	trainingParams := som.TrainingConfig{
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Linear{Start: 10, End: 0.5},
	}
	trainer, err := som.NewTrainer(&s, []*som.Table{table}, &trainingParams, rng)
	if err != nil {
		panic(err)
	}

	trainer.Train(1000)

	fmt.Println(s.Layers()[0].ToCSV(';'))
	//fmt.Println(s.Layers()[0].ColumnMatrix(0))
}

func generateData(rows, cols int) []float64 {
	data := make([]float64, cols*rows)
	for i := 0; i < rows/2; i++ {
		data[i*cols] = rand.NormFloat64()*0.1 + 0.2
		data[i*cols+1] = rand.NormFloat64()*0.1 + 0.3
	}
	for i := rows / 2; i < rows; i++ {
		data[i*cols] = rand.NormFloat64()*0.2 + 0.7
		data[i*cols+1] = rand.NormFloat64()*0.1 + 0.8
	}
	return data
}
