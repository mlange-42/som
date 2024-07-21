package som

import (
	"math"

	"github.com/mlange-42/som/table"
)

// Predictor is a struct that holds a SOM and a set of tables for making predictions.
type Predictor struct {
	som    *Som
	tables []*table.Table
}

// NewPredictor creates a new Predictor instance with the given SOM and tables.
// The tables must have the same number of rows as the SOM has nodes.
// An error is returned if the tables do not match the SOM.
func NewPredictor(som *Som, tables []*table.Table) (*Predictor, error) {
	if err := checkTables(som, tables); err != nil {
		return nil, err
	}
	return &Predictor{
		som:    som,
		tables: tables,
	}, nil
}

// Som returns the SOM associated with this Predictor.
func (p *Predictor) Som() *Som {
	return p.som
}

// GetBMUTable returns a table with the best matching units (BMUs) for each row in the
// associated tables. The table contains the following columns:
//
// - node_id: the index of the BMU node
// - node_x: the x-coordinate of the BMU node
// - node_y: the y-coordinate of the BMU node
// - node_dist: the distance between the input data and the BMU node
func (p *Predictor) GetBMUTable() *table.Table {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	cols := 4
	bmu := make([]float64, rows*cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < len(p.tables); j++ {
			data[j] = p.tables[j].GetRow(i)
		}
		idx, dist := p.som.GetBMU(data)
		x, y := p.som.Size().Coords(idx)
		bmu[i*cols] = float64(idx)
		bmu[i*cols+1] = float64(x)
		bmu[i*cols+2] = float64(y)
		bmu[i*cols+3] = dist
	}

	t, err := table.NewWithData([]string{"node_id", "node_x", "node_y", "node_dist"}, bmu)
	if err != nil {
		panic(err)
	}

	return t
}

// GetBMU returns a slice of the best matching unit (BMU) indices for each row in the
// associated tables.
func (p *Predictor) GetBMU() []int {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	bmu := make([]int, rows)

	for i := 0; i < rows; i++ {
		for j := 0; j < len(p.tables); j++ {
			data[j] = p.tables[j].GetRow(i)
		}
		idx, _ := p.som.GetBMU(data)
		bmu[i] = idx
	}

	return bmu
}

// GetBMUWithDistance returns the best matching unit (BMU) indices and the distances
// between the input data and the BMU for each row in the associated tables.
func (p *Predictor) GetBMUWithDistance() ([]int, []float64) {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	bmu := make([]int, rows)
	distance := make([]float64, rows)

	for i := 0; i < rows; i++ {
		for j := 0; j < len(p.tables); j++ {
			data[j] = p.tables[j].GetRow(i)
		}
		bmu[i], distance[i] = p.som.GetBMU(data)
	}

	return bmu, distance
}

// GetDensity returns the density of the SOM, which is the number of data points
// that map to each node in the SOM. The returned slice has one element for each
// node in the SOM, where the value at index i represents the number of data points
// that map to the node at index i.
func (p *Predictor) GetDensity() []int {
	bmu := p.GetBMU()
	counter := make([]int, p.Som().Size().Nodes())
	for _, idx := range bmu {
		counter[idx]++
	}
	return counter
}

// GetError returns the error for each node in the SOM, either as the raw sum of squared
// distances between the input data and the BMU, or as the root mean squared error
// (RMSE). The returned slice has one element for each
// node in the SOM, where the value at index i represents the error for the node
// at index i.
//
// If rmse is true, the returned values will be the RMSE .
// Otherwise, the returned values will be the MSE.
func (p *Predictor) GetError(rmse bool) []float64 {
	bmu, dist := p.GetBMUWithDistance()

	errors := make([]float64, p.som.size.Nodes())
	counter := make([]int, p.som.size.Nodes())
	for i, b := range bmu {
		d := dist[i]
		errors[b] = d * d
		counter[b]++
	}
	for i := range errors {
		if counter[i] == 0 {
			continue
		}
		errors[i] /= float64(counter[i])
		if rmse {
			errors[i] = math.Sqrt(errors[i])
		}
	}
	return errors
}
