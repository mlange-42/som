package som

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"

	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/distance"
	"github.com/mlange-42/som/layer"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/table"
)

// TrainingConfig holds the configuration parameters for training a Self-Organizing Map (SOM).
type TrainingConfig struct {
	Epochs             int         // Number of training epochs
	LearningRate       decay.Decay // Learning rate decay function
	NeighborhoodRadius decay.Decay // Neighborhood radius decay function
	WeightDecay        decay.Decay // Weight decay coefficient decay function
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
	center [][]float64
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
func (t *Trainer) Train(progress chan TrainingProgress) {
	t.som.Randomize(t.rng)

	t.calcDataCenter()

	var meanDist float64
	var qError float64
	var p TrainingProgress
	for epoch := 0; epoch < t.params.Epochs; epoch++ {
		alpha := t.params.LearningRate.Decay(epoch, t.params.Epochs)
		radius := t.params.NeighborhoodRadius.Decay(epoch, t.params.Epochs)
		decay := 0.0
		if t.params.WeightDecay != nil {
			decay = t.params.WeightDecay.Decay(epoch, t.params.Epochs)
		}

		if decay > 0 {
			t.decayWeights(decay)
		}
		meanDist, qError = t.epoch(alpha, radius)

		p.Epoch = epoch
		p.Alpha = alpha
		p.Radius = radius
		p.WeightDecay = decay
		p.MeanDist = meanDist
		p.Error = qError

		progress <- p
	}

	close(progress)
}

func (t *Trainer) calcDataCenter() {
	if t.params.WeightDecay == nil {
		return
	}

	t.center = make([][]float64, len(t.tables))
	rows := t.tables[0].Rows()

	for i, tab := range t.tables {
		cols := tab.Columns()
		t.center[i] = make([]float64, cols)
		cnt := make([]int, cols)
		for j := 0; j < rows; j++ {
			for k := 0; k < cols; k++ {
				v := tab.Get(j, k)
				if math.IsNaN(v) {
					continue
				}
				t.center[i][k] += v
				cnt[k]++
			}
		}
		for k := 0; k < cols; k++ {
			t.center[i][k] /= float64(cnt[k])
		}
	}
}

func (t *Trainer) PropagateLabels(name string, classes []string, indices []int) error {
	if len(indices) != t.tables[0].Rows() {
		return fmt.Errorf("length of indices (%d) does not match number of data rows (%d)", len(indices), t.tables[0].Rows())
	}

	probabilities, counts, err := t.findLabels(classes, indices)
	if err != nil {
		return err
	}
	lay, err := t.propagateLabels(name, probabilities, counts, classes)
	if err != nil {
		return err
	}

	t.som.layers = append(t.som.layers, lay)

	return nil
}

func (t *Trainer) propagateLabels(name string, probabilities []float64, counts []int, classes []string) (*layer.Layer, error) {
	rows := len(counts)
	cols := len(classes)
	if len(probabilities) != rows*cols {
		return nil, fmt.Errorf("length of data (%d) does not match expected length (%d*%d=%d)",
			len(probabilities), rows, cols, rows*cols)
	}

	lay1, lay2, err := t.createLabelLayers(name, classes, probabilities)
	if err != nil {
		return nil, err
	}

	uMatrix := t.som.UMatrix(false)
	sigma := t.calcPropagationSigma(uMatrix)
	neigh := &neighborhood.Gaussian{}

	w, h := t.som.Size().Width, t.som.Size().Height

	for iter := 0; iter < 10000; iter++ {
		totalDiff := 0.0
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				nodeIdx := t.som.Size().Index(x, y)
				if counts[nodeIdx] > 0 {
					// Node with known labels.
					continue
				}
				self := lay2.GetNode(x, y)
				selfPrev := lay1.GetNode(x, y)
				for i := 0; i < cols; i++ {
					totalDiff += math.Abs(self[i] - selfPrev[i])
					self[i] = 0
				}

				sumWeights := t.updateLabelsFromNeighbors(x, y, self, lay1, uMatrix, sigma, neigh)
				if sumWeights == 0 {
					continue
				}

				for i := 0; i < cols; i++ {
					self[i] /= sumWeights
				}
			}
		}

		lay1, lay2 = lay2, lay1
		if iter > 0 && totalDiff < 0.001 {
			break
		}
	}

	return lay1, nil
}

func (t *Trainer) updateLabelsFromNeighbors(x, y int,
	self []float64, lay1 *layer.Layer, uMatrix [][]float64, sigma float64,
	neigh neighborhood.Neighborhood) float64 {
	sumWeights := 0.0

	w, h := t.som.Size().Width, t.som.Size().Height
	dxMin, dxMax := max(x-1, 0)-x, min(x+1, w-1)-x
	dyMin, dyMax := max(y-1, 0)-y, min(y+1, h-1)-y
	for dx := dxMin; dx <= dxMax; dx++ {
		for dy := dyMin; dy <= dyMax; dy++ {
			if dx != 0 && dy != 0 {
				continue // diagonal neighbors
			}
			var weight float64
			if dx == 0 && dy == 0 {
				// BMU
				weight = neigh.Weight(0, sigma)
			} else {
				// Neighbor node
				weight = neigh.Weight(uMatrix[2*y+dy][2*x+dx], sigma)
			}

			other := lay1.GetNode(x+dx, y+dy)

			sumWeights += t.updateLabels(self, other, weight)
		}
	}
	return sumWeights
}

func (t *Trainer) updateLabels(self, other []float64, weight float64) float64 {
	sumWeights := 0.0
	for i := 0; i < len(self); i++ {
		v := weight * other[i]
		self[i] += v
		sumWeights += v
	}
	return sumWeights
}

func (t *Trainer) createLabelLayers(name string, classes []string, probabilities []float64) (*layer.Layer, *layer.Layer, error) {
	lay1, err := layer.NewWithData(name, classes, nil, *t.som.Size(), &distance.Hamming{}, 0.0, true, probabilities)
	if err != nil {
		return nil, nil, err
	}
	lay2, err := layer.NewWithData(name, classes, nil, *t.som.Size(), &distance.Hamming{}, 0.0, true, append(make([]float64, 0, len(probabilities)), probabilities...))
	if err != nil {
		return nil, nil, err
	}
	return lay1, lay2, nil
}

func (t *Trainer) calcPropagationSigma(uMatrix [][]float64) float64 {
	w, h := t.som.Size().Width, t.som.Size().Height
	values := make([]float64, (w-1)*h+w*(h-1))
	idx := 0
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if x < w-1 {
				values[idx] = uMatrix[y*2][x*2+1]
				idx++
			}
			if y < h-1 {
				values[idx] = uMatrix[y*2+1][x*2]
				idx++
			}
		}
	}
	slices.Sort(values)

	sMin := values[0]
	sMax := values[len(values)-1]

	urpMin := math.Inf(1)
	urpMinIndex := -1
	for i := range values {
		urp := calcSquaredURP(sMin, sMax, values, i)
		if urp < urpMin {
			urpMin = urp
			urpMinIndex = i
		}
	}

	return values[urpMinIndex] / 3.0
}

func calcSquaredURP(sMin, sMax float64, distribution []float64, index int) float64 {
	sigma := distribution[index]
	v1 := (sigma - sMin) / (sMax - sMin)
	v2 := 1 - (float64(index) / float64(len(distribution)))

	d2 := v1*v1 + v2*v2
	return d2
}

func (t *Trainer) findLabels(classes []string, indices []int) ([]float64, []int, error) {
	pred, err := NewPredictor(t.som, t.tables)
	if err != nil {
		return nil, nil, err
	}

	bmu := pred.GetBMU()

	rows := t.som.Size().Nodes()
	cols := len(classes)

	classCounter := make([]float64, rows*cols)
	totalCounter := make([]int, rows)
	for i, v := range bmu {
		classCounter[v*cols+indices[i]]++
		totalCounter[v]++
	}
	for i := 0; i < rows; i++ {
		if totalCounter[i] == 0 {
			continue
		}
		for j := 0; j < cols; j++ {
			classCounter[i*cols+j] /= float64(totalCounter[i])
		}
	}

	return classCounter, totalCounter, nil
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

func (t *Trainer) decayWeights(beta float64) {
	t.som.decayWeights(t.center, beta)
}

// TrainingProgress represents the progress of a training epoch.
type TrainingProgress struct {
	Epoch       int     // The current epoch number
	Alpha       float64 // The current learning rate alpha
	Radius      float64 // The current neighborhood radius
	WeightDecay float64 // The weight decay factor
	MeanDist    float64 // The mean distance of the training data to the SOM
	Error       float64 // The quantization error (MSE)
}

// CsvHeader returns a CSV header row for the TrainingProgress struct fields, using the provided delimiter.
func (p *TrainingProgress) CsvHeader(delim rune) string {
	return fmt.Sprintf("Epoch%cAlpha%cRadius%cDecay%cMeanDist%cError", delim, delim, delim, delim, delim)
}

// CsvRow returns a comma-separated string representation of the TrainingProgress struct fields.
// The values are formatted using the provided delimiter character.
func (p *TrainingProgress) CsvRow(delim rune) string {
	return fmt.Sprintf("%d%c%s%c%s%c%s%c%s%c%s",
		p.Epoch, delim,
		strconv.FormatFloat(p.Alpha, 'f', -1, 64), delim,
		strconv.FormatFloat(p.Radius, 'f', -1, 64), delim,
		strconv.FormatFloat(p.WeightDecay, 'f', -1, 64), delim,
		strconv.FormatFloat(p.MeanDist, 'f', -1, 64), delim,
		strconv.FormatFloat(p.Error, 'f', -1, 64))
}
