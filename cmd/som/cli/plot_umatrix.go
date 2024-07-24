package cli

import (
	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func plotUMatrixCommand() *cobra.Command {
	var size []int
	var dataFile string
	var labelsColumn string
	var delim string
	var noData string
	var ignore []string
	var sample int

	command := &cobra.Command{
		Use:   "u-matrix [flags] <som-file> <out-file>",
		Short: "Plots the u-matrix of an SOM, showing inter-nde distances",
		Long:  `Plots the u-matrix of an SOM, showing inter-nde distances`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			outFile := args[1]

			return plotHeatmap(size,
				somFile, outFile, dataFile,
				labelsColumn, delim, noData, "U-Matrix",
				ignore, sample,
				func(s *som.Som, p *som.Predictor, r table.Reader) (plotter.GridXYZ, []string, error) {
					uMatrix := s.UMatrix(true)
					return &plot.UMatrixGrid{UMatrix: uMatrix}, nil, nil
				},
			)
		},
	}

	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of the plot in pixels")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required for --labels")
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")
	command.Flags().StringVarP(&labelsColumn, "labels", "l", "", "Labels column in the data file")
	command.Flags().IntVarP(&sample, "sample", "S", 0, "Sample this many rows from the data file (default all)")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false
	command.MarkFlagFilename("data-file", "csv")

	return command
}
