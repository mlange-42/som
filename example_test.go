package som_test

import (
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"os"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/norm"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
)

// Trains an SOM from randomly generated 2D data,
// and plots a scatter plot of the data and the trained SOM.
func Example() {
	// Create a random number generator for reproductible results
	rng := rand.New(rand.NewSource(42))

	// Generate random data for training
	data := generateRandomData("a", "b", 250, rng)

	// Create an SOM configuration, matching table columns.
	config := som.SomConfig{
		Size: layer.Size{
			Width:  12,
			Height: 8,
		},
		Neighborhood: &neighborhood.Gaussian{},
		MapMetric:    &neighborhood.ManhattanMetric{},
		Layers: []*som.LayerDef{
			{
				Name: "L1",
				Columns: []string{
					"a",
					"b",
				},
				Norm: []norm.Normalizer{
					&norm.Identity{},
					&norm.Identity{},
				},
				Metric: &distance.Euclidean{},
			},
		},
	}

	// Create an SOM from the configuration
	s, err := som.New(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a training configuration
	trainingConfig := som.TrainingConfig{
		Epochs:             1000,
		LearningRate:       &decay.Polynomial{Start: 0.5, End: 0.01, Exp: 2},
		NeighborhoodRadius: &decay.Polynomial{Start: 6, End: 0.5, Exp: 2},
		ViSomLambda:        0.0,
	}

	// Create a trainer instance from SOM and training data
	trainer, err := som.NewTrainer(s, []*table.Table{data}, &trainingConfig, rng)
	if err != nil {
		log.Fatal(err)
	}
	_ = trainer

	// Create a channel for training progress updates
	progress := make(chan som.TrainingProgress)

	// Run SOM training asynchronously
	go trainer.Train(progress)

	// Wait for training to finish
	for p := range progress {
		if p.Epoch%100 == 0 {
			fmt.Printf("Epoch %03d (err=%.4f)\n", p.Epoch, p.Error)
		}
	}

	// Create data sources for plotting
	xy := plot.SomXY{Som: s, XLayer: 0, XColumn: 0, YLayer: 0, YColumn: 1}
	dataXY := plot.TableXY{XTable: data, YTable: data, XColumn: 0, YColumn: 1, XNorm: &norm.Identity{}, YNorm: &norm.Identity{}}
	// Make a scatter plot
	img, err := plot.XY("xy", &xy, *s.Size(), 600, 400, nil, nil, true, &dataXY, nil, nil, false)
	if err != nil {
		log.Fatal(err)
	}

	// Open a file to write the image
	file, err := os.Create("som-xy.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write the image to a PNG file
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	//Epoch 000 (err=0.1164)
	//Epoch 100 (err=0.0641)
	//Epoch 200 (err=0.0370)
	//Epoch 300 (err=0.0206)
	//Epoch 400 (err=0.0119)
	//Epoch 500 (err=0.0064)
	//Epoch 600 (err=0.0036)
	//Epoch 700 (err=0.0022)
	//Epoch 800 (err=0.0015)
	//Epoch 900 (err=0.0012)
}

// Generate 2D random data for training.
func generateRandomData(xCol, yCol string, rows int, rng *rand.Rand) *table.Table {
	data := make([]float64, rows*2)
	for i := 0; i < rows; i++ {
		x := rng.Float64()*2 - 1
		y := x*x + rng.NormFloat64()*0.1
		data[i*2] = x
		data[i*2+1] = y
	}
	t, err := table.NewWithData([]string{xCol, yCol}, data)
	if err != nil {
		panic(err)
	}
	return t
}
