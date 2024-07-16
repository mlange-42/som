package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/neighborhood"
)

func main() {
	rng := rand.New(rand.NewSource(1))

	somParams := som.SomParams{
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

	trainingParams := som.TrainingParams{
		LearningRate:       &decay.Linear{Start: 0.5, End: 0.01},
		NeighborhoodRadius: &decay.Linear{Start: 10, End: 0.5},
	}
	trainer, err := som.NewTrainer(&s, &table, &trainingParams, rng)
	if err != nil {
		panic(err)
	}

	//fmt.Println(printTable(&table))

	trainer.Train(10)

	fmt.Println(printLayer(&s.Layers()[0], 0, 1))
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

func printTable(t *som.Table) string {
	b := strings.Builder{}
	for _, col := range t.Columns() {
		b.WriteString(fmt.Sprintf("%12s", col))
	}
	b.WriteRune('\n')
	cols := len(t.Columns())
	for i := 0; i < t.Rows(); i++ {
		for j := 0; j < cols; j++ {
			b.WriteString(fmt.Sprintf("%12.2f", t.Get(i, j)))
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func printLayer(l *som.Layer, x, y int) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%12s", l.Columns()[x]))
	b.WriteString(fmt.Sprintf("%12s", l.Columns()[y]))
	b.WriteRune('\n')

	nodes := l.Nodes()
	for i := 0; i < nodes; i++ {
		b.WriteString(fmt.Sprintf("%12.2f", l.GetAt(i, x)))
		b.WriteString(fmt.Sprintf("%12.2f", l.GetAt(i, y)))
		b.WriteRune('\n')
	}
	return b.String()
}
