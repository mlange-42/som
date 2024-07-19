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

func uMatrixCommand() *cobra.Command {
	var size []int
	var dataFile string
	var labelsColumn string
	var delim string
	var noData string

	command := &cobra.Command{
		Use:   "u-matrix [flags] <som-file> <out-file>",
		Short: "Plots the u-matrix of an SOM",
		Long:  `Plots the u-matrix of an SOM`,
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

			if dataFile != "" {
				var err error
				reader, err = csv.NewFileReader(dataFile, del[0], noData)
				if err != nil {
					return err
				}
				predictor, _, err = createPredictor(config, s, reader)
				if err != nil {
					return err
				}
			}

			var labels []string
			var positions []plotter.XY

			if labelsColumn != "" {
				if dataFile == "" {
					return fmt.Errorf("data file must be specified when labels column is specified")
				}

				labels, positions, err = extractLabels(predictor, labelsColumn, reader)
				if err != nil {
					return err
				}
			}

			uMatrix := s.UMatrix()
			grid := &plot.UMatrixGrid{UMatrix: uMatrix}
			title := "U-Matrix"

			img, err := plot.Heatmap(title, grid, size[0], size[1], nil, labels, positions)
			if err != nil {
				return err
			}

			return writeImage(img, outFile)
		},
	}

	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required for --labels")
	command.Flags().StringVarP(&labelsColumn, "labels", "l", "", "Labels column in the data file")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false

	return command
}
