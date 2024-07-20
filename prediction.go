package som

import "github.com/mlange-42/som/table"

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

	cols := 3
	bmu := make([]float64, rows*cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < len(p.tables); j++ {
			data[j] = p.tables[j].GetRow(i)
		}
		idx, _ := p.som.getBMU(data)
		x, y := p.som.Size().Coords(idx)
		bmu[i*cols] = float64(idx)
		bmu[i*cols+1] = float64(x)
		bmu[i*cols+2] = float64(y)
	}

	return table.NewWithData([]string{"node_id", "node_x", "node_y"}, bmu)
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
