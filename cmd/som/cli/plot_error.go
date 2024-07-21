package cli

import (
	"fmt"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func errorCommand() *cobra.Command {
	var size []int
	var dataFile string
	var labelsColumn string
	var delim string
	var noData string
	var rmse bool

	command := &cobra.Command{
		Use:   "error [flags] <som-file> <out-file>",
		Short: "Plots mean-squared node error as a heatmap",
		Long:  `Plots mean-squared node error as a heatmap`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			outFile := args[1]

			del := []rune(delim)
			if len(delim) != 1 {
				return fmt.Errorf("delimiter must be a single character")
			}
			if len(size) != 2 {
				return fmt.Errorf("size must be two integers")
			}

			config, s, err := readSom(somFile)
			if err != nil {
				return err
			}

			var reader table.Reader
			var predictor *som.Predictor

			reader, err = csv.NewFileReader(dataFile, del[0], noData)
			if err != nil {
				return err
			}
			predictor, _, err = createPredictor(config, s, reader)
			if err != nil {
				return err
			}

			var labels []string
			var positions []plotter.XY

			if labelsColumn != "" {
				labels, positions, err = extractLabels(predictor, labelsColumn, reader)
				if err != nil {
					return err
				}
			}

			mse := predictor.GetError(rmse)
			grid := &plot.FloatGrid{Size: *s.Size(), Values: mse}
			title := "Mean squared error"

			img, err := plot.Heatmap(title, grid, size[0], size[1], nil, labels, positions)
			if err != nil {
				return err
			}

			return writeImage(img, outFile)
		},
	}

	command.Flags().BoolVarP(&rmse, "rmse", "r", false, "Use root mean squared error instead of mean squared error")

	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required")
	command.Flags().StringVarP(&labelsColumn, "labels", "l", "", "Labels column in the data file")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false
	command.MarkFlagRequired("data-file")
	command.MarkFlagFilename("data-file", "csv")

	return command
}
