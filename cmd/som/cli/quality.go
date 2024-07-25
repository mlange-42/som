package cli

import (
	"fmt"
	"os"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/neighborhood"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func qualityCommand() *cobra.Command {
	var delim string
	var noData string
	var ignore []string

	command := &cobra.Command{
		Use:   "quality [flags] <som-file> <data-file>",
		Short: "Calculates various quality metrics for a trained SOM.",
		Long: `Calculates various quality metrics for a trained SOM.

Calculates the following SOM quality metrics and prints them to STDOUT:

 - Quantization error
 - Mean square error
 - Root mean square error
 - Topographic error`,
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
			eval := som.NewEvaluator(pred)

			qe, mse, rmse := eval.Error()
			te := eval.TopographicError(&neighborhood.ManhattanMetric{})

			fmt.Printf(`Quantization error:     %7.3f
Mean square error:      %7.3f
Root mean square error: %7.3f
Topographic error:      %7.3f
`, qe, mse, rmse, te)

			return nil
		},
	}
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter for CSV input and output")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data string for CSV input and output")

	command.Flags().SortFlags = false

	return command
}
