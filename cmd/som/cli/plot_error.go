package cli

import (
	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
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

			return plotHeatmap(size,
				somFile, outFile, dataFile,
				labelsColumn, delim, noData, "U-Matrix",
				func(s *som.Som, p *som.Predictor) plotter.GridXYZ {
					mse := p.GetError(rmse)
					return &plot.FloatGrid{Size: *s.Size(), Values: mse}
				},
			)
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
