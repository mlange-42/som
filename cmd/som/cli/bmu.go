package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/table"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func bmuCommand() *cobra.Command {
	var delim string
	var noData string
	var preserve []string
	var ignore []string

	command := &cobra.Command{
		Use:   "bmu [flags] <som-file> <data-file>",
		Short: "Finds the best-matching unit (BMU) for each table row in a dataset.",
		Long: `Finds the best-matching unit (BMU) for each table row in a dataset.

A table with the same row order as the input table is created. Columns are:

 - node_id: the index of the BMU node
 - node_x: the x-coordinate of the BMU node
 - node_y: the y-coordinate of the BMU node
 - node_dist: the distance between the input data and the BMU node

The result table is written to STDOUT in CSV format.
Redirect output to a file like this:
 
  som bmu som.yml data.csv > bmu.csv

Columns from the input table can be transferred to the output table by use of
the --preserve flag. Here is how to transfer 'ID' and 'Name' columns:

  sum bmu som.yml data.csv --preserve ID,Name > bmu.csv
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			dataFile := args[1]

			somYaml, err := os.ReadFile(somFile)
			if err != nil {
				return err
			}
			config, _, err := yml.ToSomConfig(somYaml)
			if err != nil {
				return err
			}

			del := []rune(delim)
			if len(delim) != 1 {
				return fmt.Errorf("delimiter must be a single character")
			}

			reader, err := csv.NewFileReader(dataFile, del[0], noData)
			if err != nil {
				return err
			}

			preserved := [][]string{}
			for _, column := range preserve {
				col, err := reader.ReadLabels(column)
				if err != nil {
					return err
				}
				preserved = append(preserved, col)
			}

			tables, _, err := config.PrepareTables(reader, ignore, false, false)
			if err != nil {
				return err
			}

			s, err := som.New(config)
			if err != nil {
				return err
			}
			pred, err := som.NewPredictor(s, tables)
			if err != nil {
				return err
			}

			bmu := pred.GetBMUTable()
			writer := strings.Builder{}
			err = csv.TablesToCsv([]*table.Table{bmu}, preserve, preserved, &writer, del[0], noData)
			if err != nil {
				return err
			}

			fmt.Println(writer.String())

			return nil
		},
	}
	command.Flags().StringSliceVarP(&preserve, "preserve", "p", nil, "Preserve columns and prepend them to the output table")
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter for CSV input and output")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data string for CSV input and output")

	command.Flags().SortFlags = false

	return command
}
