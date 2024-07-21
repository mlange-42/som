package som

import (
	"math"

	"github.com/mlange-42/som/table"
)

type Predictor struct {
	som    *Som
	tables []*table.Table
}

func NewPredictor(som *Som, tables []*table.Table) (*Predictor, error) {
	if err := checkTables(som, tables); err != nil {
		return nil, err
	}
	return &Predictor{
		som:    som,
		tables: tables,
	}, nil
}

func (p *Predictor) Som() *Som {
	return p.som
}

func (p *Predictor) GetBMUTable() (*table.Table, error) {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	cols := 4
	bmu := make([]float64, rows*cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < len(p.tables); j++ {
			data[j] = p.tables[j].GetRow(i)
		}
		idx, dist := p.som.getBMU(data)
		x, y := p.som.Size().Coords(idx)
		bmu[i*cols] = float64(idx)
		bmu[i*cols+1] = float64(x)
		bmu[i*cols+2] = float64(y)
		bmu[i*cols+3] = dist
	}

	return table.NewWithData([]string{"node_id", "node_x", "node_y", "node_dist"}, bmu)
}

func (p *Predictor) GetBMU() []int {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	bmu := make([]int, rows)

	for i := 0; i < rows; i++ {
		for j := 0; j < len(p.tables); j++ {
			data[j] = p.tables[j].GetRow(i)
		}
		idx, _ := p.som.getBMU(data)
		bmu[i] = idx
	}

	return bmu
}

func (p *Predictor) GetBMUDistance() ([]int, []float64) {
	data := make([][]float64, len(p.tables))
	rows := p.tables[0].Rows()

	bmu := make([]int, rows)
	distance := make([]float64, rows)

	for i := 0; i < rows; i++ {
		for j := 0; j < len(p.tables); j++ {
			data[j] = p.tables[j].GetRow(i)
		}
		bmu[i], distance[i] = p.som.getBMU(data)
	}

	return bmu, distance
}

func (p *Predictor) GetDensity() []int {
	bmu := p.GetBMU()
	counter := make([]int, p.Som().Size().Nodes())
	for _, idx := range bmu {
		counter[idx]++
	}
	return counter
}

func (p *Predictor) GetError(rmse bool) []float64 {
	bmu, dist := p.GetBMUDistance()

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
