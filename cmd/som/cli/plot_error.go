package cli

import (
	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func plotErrorCommand() *cobra.Command {
	var size []int
	var dataFile string
	var labelsColumn string
	var boundaries string
	var delim string
	var noData string
	var rmse bool
	var ignore []string
	var sample int

	command := &cobra.Command{
		Use:   "error [flags] <som-file> <out-file>",
		Short: "Plots (root) mean-squared node error as a heatmap.",
		Long: `Plots (root) mean-squared node error as a heatmap.

By default, the mean-squared error is plotted.
With --rmse, the root mean squared error is shown instead.

Data provided via --data-file can be displayed on top of the error plot,
showing the values in the column given by the --label flag:

  som plot error som.yml error.png --data-file data.csv --label name

For large datasets, --sample can be used to show only a sub-set of the data.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			outFile := args[1]

			var title string
			if rmse {
				title = "Root Mean Squared Error"
			} else {
				title = "Mean Squared Error"
			}

			return plotHeatmap(size,
				somFile, outFile, dataFile,
				labelsColumn, delim, noData, title,
				ignore, boundaries, sample,
				func(s *som.Som, p *som.Predictor, r table.Reader) (plotter.GridXYZ, []string, error) {
					mse := p.GetError(rmse)
					return &plot.FloatGrid{Size: *s.Size(), Values: mse}, nil, nil
				},
			)
		},
	}

	command.Flags().BoolVarP(&rmse, "rmse", "r", false, "Use root mean squared error instead of mean squared error")
	command.Flags().StringVarP(&boundaries, "boundaries", "b", "", "Optional categorical variable to show boundaries for")

	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of the plot in pixels")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required")
	command.Flags().StringVarP(&labelsColumn, "label", "l", "", "Label column in the data file")
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")
	command.Flags().IntVarP(&sample, "sample", "S", 0, "Sample this many rows from the data file (default all)")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false
	command.MarkFlagRequired("data-file")
	command.MarkFlagFilename("data-file", "csv")

	return command
}
