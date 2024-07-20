package cli

import (
	"fmt"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func xyCommand() *cobra.Command {
	var size []int
	var xColumn string
	var yColumn string
	var color string
	var dataFile string
	var labelsColumn string
	var delim string
	var noData string

	command := &cobra.Command{
		Use:   "xy [flags] <som-file> <out-file>",
		Short: "Scatter plots for pairs of SOM variables",
		Long:  `Scatter plots for pairs of SOM variables`,
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

			columns := []string{xColumn, yColumn}
			if color != "" {
				columns = append(columns, color)
			}

			_, indices, err := extractIndices(s, columns)
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

			title := fmt.Sprintf("%s vs. %s", columns[0], columns[1])
			xy := plot.SomXY{
				Som:     s,
				XLayer:  indices[0][0],
				XColumn: indices[0][1],
				YLayer:  indices[1][0],
				YColumn: indices[1][1],
			}

			var classes []string
			var classIndices []int
			if len(indices) > 2 {
				classes, classIndices = conv.LayerToClasses(s.Layers()[indices[2][0]])
			}

			img, err := plot.XY(title, &xy, size[0], size[1], classes, classIndices, labels, positions)
			if err != nil {
				return err
			}

			return writeImage(img, outFile)
		},
	}

	command.Flags().StringVarP(&xColumn, "x-column", "x", "x", "Column for x axis")
	command.Flags().StringVarP(&yColumn, "y-column", "y", "y", "Column for y axis")
	command.Flags().StringVarP(&color, "color", "c", "", "Column for color")

	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required for --labels")
	command.Flags().StringVarP(&labelsColumn, "labels", "l", "", "Labels column in the data file")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false

	return command
}
