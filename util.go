package som

import (
	"fmt"

	"github.com/mlange-42/som/table"
)

// checkTable checks that the table columns match the SOM layer columns.
// It returns true if the table columns match the SOM layer columns, false otherwise.
func checkTables(som *Som, tables []*table.Table) error {
	if len(tables) == 0 {
		return fmt.Errorf("no tables provided")
	}

	if len(som.layers) != len(tables) {
		return fmt.Errorf("number of tables (%d) does not match number of layers (%d)", len(tables), len(som.layers))
	}

	rows := -1
	for _, table := range tables {
		if table == nil {
			continue
		}
		if rows == -1 {
			rows = table.Rows()
		} else if rows != table.Rows() {
			return fmt.Errorf("number of rows in table (%d) does not match number of rows in table (%d)", rows, table.Rows())
		}
	}

	for i := range som.layers {
		table := tables[i]
		if table == nil {
			continue
		}
		cols := som.layers[i].ColumnNames()
		if table.Columns() != len(cols) {
			return fmt.Errorf("number of columns in table (%d) does not match number of columns in layer (%d)", table.Columns(), len(cols))
		}
		for j, col := range cols {
			if table.ColumnNames()[j] != col {
				return fmt.Errorf("column %d in table (%s) does not match column %d in layer (%s)", j, table.ColumnNames()[j], j, col)
			}
		}
	}
	return nil
}
