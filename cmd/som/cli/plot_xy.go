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
	var dataColor string
	var noGrid bool
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
			var tables []*table.Table
			var data plotter.XYer
			var dataCats []string
			var dataIndices []int

			if dataFile != "" {
				var err error
				reader, err = csv.NewFileReader(dataFile, del[0], noData)
				if err != nil {
					return err
				}
				tables, err = config.PrepareTables(reader, false)
				if err != nil {
					return err
				}
				data = extractData(config, tables, indices)

				if dataColor != "" {
					colorColumn, err := reader.ReadLabels(dataColor)
					if err != nil {
						return err
					}
					dataCats, dataIndices = conv.ClassesToIndices(colorColumn)
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

			img, err := plot.XY(
				title, &xy, *s.Size(), size[0], size[1],
				classes, classIndices, !noGrid,
				data, dataCats, dataIndices, color != dataColor)
			if err != nil {
				return err
			}

			return writeImage(img, outFile)
		},
	}

	command.Flags().StringVarP(&xColumn, "x-column", "x", "x", "Column for x axis")
	command.Flags().StringVarP(&yColumn, "y-column", "y", "y", "Column for y axis")
	command.Flags().StringVarP(&color, "color", "c", "", "Column for color")

	command.Flags().BoolVarP(&noGrid, "no-grid", "G", false, "Don't draw SOM grid lines")
	command.Flags().StringVarP(&dataColor, "data-color", "C", "", "Column for data color")

	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required for --labels")
	command.Flags().StringVarP(&labelsColumn, "labels", "l", "", "Labels column in the data file")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false

	return command
}

func extractData(config *som.SomConfig, tables []*table.Table, indices [][2]int) plotter.XYer {
	return &plot.TableXY{
		XTable:  tables[indices[0][0]],
		YTable:  tables[indices[1][0]],
		XColumn: indices[0][1],
		YColumn: indices[1][1],
		XNorm:   config.Layers[indices[0][0]].Norm[indices[0][1]],
		YNorm:   config.Layers[indices[1][0]].Norm[indices[1][1]],
	}
}
