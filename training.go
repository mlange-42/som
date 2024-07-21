package som

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/table"
)

// TrainingConfig holds the configuration parameters for training a Self-Organizing Map (SOM).
type TrainingConfig struct {
	Epochs             int         // Number of training epochs
	LearningRate       decay.Decay // Learning rate decay function
	NeighborhoodRadius decay.Decay // Neighborhood radius decay function
	ViSomLambda        float64     // ViSOM lambda resolution parameter
}

// Trainer is a struct that holds the necessary components for training a Self-Organizing Map (SOM).
// It contains a reference to the SOM, the training data tables, the training configuration parameters,
// and a random number generator.
type Trainer struct {
	som    *Som
	tables []*table.Table
	params *TrainingConfig
	rng    *rand.Rand
}

// NewTrainer creates a new Trainer instance with the provided SOM, data tables, training configuration, and random number generator.
// It performs a check on the provided data tables to ensure they are compatible with the SOM.
// If the check fails, an error is returned.
func NewTrainer(som *Som, tables []*table.Table, params *TrainingConfig, rng *rand.Rand) (*Trainer, error) {
	if err := checkTables(som, tables); err != nil {
		return nil, err
	}
	return &Trainer{
		som:    som,
		tables: tables,
		params: params,
		rng:    rng,
	}, nil
}

// Train trains the Self-Organizing Map (SOM) using the provided training data and configuration.
// It iterates through the specified number of epochs, updating the learning rate and neighborhood radius
// at each epoch. For each epoch, it performs a single training iteration,
// and sends the training progress information (epoch, learning rate, neighborhood radius, mean distance,
// and quantization error) to the provided progress channel.
// After all epochs are completed, the channel is closed.
func (t *Trainer) Train(progress chan *TrainingProgress) {
	t.som.Randomize(t.rng)

	var meanDist float64
	var qError float64
	var p TrainingProgress
	for epoch := 0; epoch < t.params.Epochs; epoch++ {
		alpha := t.params.LearningRate.Decay(epoch, t.params.Epochs)
		radius := t.params.NeighborhoodRadius.Decay(epoch, t.params.Epochs)
		meanDist, qError = t.epoch(alpha, radius)

		p.Epoch = epoch
		p.Alpha = alpha
		p.Radius = radius
		p.MeanDist = meanDist
		p.Error = qError

		progress <- &p
	}

	close(progress)
}

func (t *Trainer) epoch(alpha, radius float64) (meanDist, quantError float64) {
	data := make([][]float64, len(t.tables))
	rows := t.tables[0].Rows()

	sumDist := 0.0
	sumDistSq := 0.0
	for i := 0; i < rows; i++ {
		for j := 0; j < len(t.tables); j++ {
			data[j] = t.tables[j].GetRow(i)
		}
		dist := t.som.Learn(data, alpha, radius, t.params.ViSomLambda)
		sumDist += dist
		sumDistSq += dist * dist

		if t.params.ViSomLambda == 0 || i%10 != 0 { // SOM
			continue
		}
		// ViSOM refresh: present random node as data
		node := t.rng.Intn(t.som.size.Nodes())
		for j := 0; j < len(t.tables); j++ {
			data[j] = t.som.layers[j].GetNodeAt(node)
		}
		t.som.Learn(data, alpha, radius, t.params.ViSomLambda)
	}

	return sumDist / float64(rows), sumDistSq / float64(rows)
}

// TrainingProgress represents the progress of a training epoch.
type TrainingProgress struct {
	Epoch    int     // The current epoch number
	Alpha    float64 // The current learning rate alpha
	Radius   float64 // The current neighborhood radius
	MeanDist float64 // The mean distance of the training data to the SOM
	Error    float64 // The quantization error (MSE)
}

// CsvHeader returns a CSV header row for the TrainingProgress struct fields, using the provided delimiter.
func (p *TrainingProgress) CsvHeader(delim rune) string {
	return fmt.Sprintf("Epoch%cAlpha%cRadius%cMeanDist%cError", delim, delim, delim, delim)
}

// CsvRow returns a comma-separated string representation of the TrainingProgress struct fields.
// The values are formatted using the provided delimiter character.
func (p *TrainingProgress) CsvRow(delim rune) string {
	return fmt.Sprintf("%d%c%s%c%s%c%s%c%s",
		p.Epoch, delim,
		strconv.FormatFloat(p.Alpha, 'f', -1, 64), delim,
		strconv.FormatFloat(p.Radius, 'f', -1, 64), delim,
		strconv.FormatFloat(p.MeanDist, 'f', -1, 64), delim,
		strconv.FormatFloat(p.Error, 'f', -1, 64))
}
