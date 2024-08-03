package som

import (
	"fmt"
	"math"
	"slices"

	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/table"
)

// Predictor is a struct that holds an SOM and a set of tables for making predictions.
type Predictor struct {
	som    *Som
	tables []*table.Table
}

// NewPredictor creates a new Predictor instance with the given SOM and tables.
// The tables must have the same number of rows as the SOM has nodes.
// Tables are assumed to be normalized.
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

// Tables returns the tables associated with this Predictor.
func (p *Predictor) Tables() []*table.Table {
	return p.tables
}

func (p *Predictor) bmu(data [][]float64) (int, float64) {
	return p.som.GetBMU(data)
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
		p.collectData(i, data)

		idx, dist := p.bmu(data)
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

// FillMissing fills in any missing values in the input tables by using the best matching units
// (BMUs) from the SOM to determine the appropriate values to fill in.
// Tables in the argument should not be normalized.
// The number of rows in the input tables must match the number of
// rows in the Predictor's tables.
func (p *Predictor) FillMissing(tables []*table.Table) error {
	if err := checkTables(p.som, tables); err != nil {
		return err
	}

	rows := tables[0].Rows()
	if rows != p.tables[0].Rows() {
		return fmt.Errorf("number of rows in tables does not match number of rows in predictor tables")
	}

	hasMissing := findRowsWithMissing(tables)

	data := make([][]float64, len(tables))
	for i := 0; i < rows; i++ {
		if !hasMissing[i] {
			continue
		}

		p.collectData(i, data)
		bmu, _ := p.bmu(data)
		for j, t := range tables {
			if t == nil {
				continue
			}
			lay := p.som.layers[j]
			node := lay.GetNodeAt(bmu)
			outRow := t.GetRow(i)
			for k := range t.Columns() {
				if math.IsNaN(outRow[k]) {
					norm := lay.Normalizers()[k]
					outRow[k] = norm.DeNormalize(node[k])
				}
			}
		}
	}

	return nil
}

// Predict generates predictions for the specified layers in the input tables using the
// self-organizing map (SOM) associated with the Predictor. The input tables should not
// be normalized. The function will create new tables for the predicted layers and
// populate them with the predicted values.
//
// If any of the layers to predict are already present in the input tables, an error
// will be returned. The number of rows in the input tables must match the number of
// rows in the Predictor's tables.
func (p *Predictor) Predict(tables []*table.Table, layers []string) error {
	if err := checkTables(p.som, tables); err != nil {
		return err
	}

	rows := tables[0].Rows()
	if rows != p.tables[0].Rows() {
		return fmt.Errorf("number of rows in tables does not match number of rows in predictor tables")
	}

	toPredict := make([]bool, len(p.som.layers))
	for i, l := range p.som.layers {
		toPredict[i] = slices.Contains(layers, l.Name())
		if !toPredict[i] {
			continue
		}
		if tables[i] != nil {
			return fmt.Errorf("layer %s to predict is already present in input", l.Name())
		}
		tables[i] = table.New(l.ColumnNames(), rows)
	}

	data := make([][]float64, len(tables))
	for i := 0; i < rows; i++ {
		p.collectData(i, data)
		bmu, _ := p.bmu(data)

		for j, lay := range p.som.layers {
			if !toPredict[j] {
				continue
			}
			tab := tables[j]
			node := lay.GetNodeAt(bmu)
			outRow := tab.GetRow(i)

			for k := range tab.Columns() {
				norm := lay.Normalizers()[k]
				outRow[k] = norm.DeNormalize(node[k])
			}
		}
	}

	return nil
}

func findRowsWithMissing(tables []*table.Table) []bool {
	rows := tables[0].Rows()
	hasMissing := make([]bool, rows)
	for _, t := range tables {
		if t == nil {
			continue
		}

		for j := 0; j < rows; j++ {
			if hasMissing[j] {
				continue
			}

			row := t.GetRow(j)
			for k := range t.Columns() {
				if math.IsNaN(row[k]) {
					hasMissing[j] = true
					break
				}
			}
		}
	}
	return hasMissing
}

func (p *Predictor) collectData(row int, data [][]float64) {
	for j := 0; j < len(p.tables); j++ {
		t := p.tables[j]
		if t == nil {
			data[j] = nil
		} else {
			data[j] = t.GetRow(row)
		}
	}
}

// GetRowBMU returns the best matching unit (BMU) index and the distance between the
// input data and the BMU for the given row in the associated tables.
func (p *Predictor) GetRowBMU(row int) (int, float64) {
	data := make([][]float64, len(p.tables))
	p.collectData(row, data)
	return p.bmu(data)
}

// GetBMU returns a slice of the best matching unit (BMU) indices for each row in the
// associated tables.
func (p *Predictor) GetBMU() []int {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	bmu := make([]int, rows)

	for i := 0; i < rows; i++ {
		p.collectData(i, data)

		idx, _ := p.bmu(data)
		bmu[i] = idx
	}

	return bmu
}

func (p *Predictor) getBMU2() []bmu2 {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	bmu := make([]bmu2, rows)

	for i := 0; i < rows; i++ {
		p.collectData(i, data)

		idx, d, idx2, d2 := p.som.GetBMU2(data)
		bmu[i] = bmu2{idx, d, idx2, d2}
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
		p.collectData(i, data)

		bmu[i], distance[i] = p.bmu(data)
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

type bmu2 struct {
	Idx1  int
	Dist1 float64
	Idx2  int
	Dist2 float64
}

type Evaluator struct {
	predictor *Predictor
	bmu       []bmu2
}

func NewEvaluator(predictor *Predictor) *Evaluator {
	return &Evaluator{
		predictor: predictor,
		bmu:       predictor.getBMU2(),
	}
}

func (e *Evaluator) Error() (qe, mse, rmse float64) {
	sumDist := 0.0
	errorSum := 0.0

	for _, b := range e.bmu {
		sumDist += b.Dist1
		errorSum += b.Dist1 * b.Dist1
	}

	return sumDist / float64(len(e.bmu)),
		errorSum / float64(len(e.bmu)),
		math.Sqrt(errorSum / float64(len(e.bmu)))
}

func (e *Evaluator) TopographicError(dist neighborhood.Metric) float64 {
	failed := len(e.bmu)
	for _, b := range e.bmu {
		x1, y1 := e.predictor.som.Size().Coords(b.Idx1)
		x2, y2 := e.predictor.som.Size().Coords(b.Idx2)

		if dist.Distance(x1, y1, x2, y2) > 1 {
			continue
		}
		failed--
	}

	return float64(failed) / float64(len(e.bmu))
}
