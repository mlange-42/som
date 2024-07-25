package cli

import (
	"fmt"
	"math"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func plotXyCommand() *cobra.Command {
	var size []int
	var xColumn string
	var yColumn string
	var color string
	var dataColor string
	var noGrid bool
	var dataFile string
	var delim string
	var noData string
	var ignore []string

	command := &cobra.Command{
		Use:   "xy [flags] <som-file> <out-file>",
		Short: "Plots for pairs of SOM variables as scatter plots.",
		Long: `Plots for pairs of SOM variables as scatter plots.

Nodes are plotted according to their values in two SOM variables.
An optional --color column can be used to color-code the nodes:

  som plot xy som.yml xy.png -x X -y Y -c Class

By default, the 2-dimensional lattice/grid of nodes is shown as lines
between neighboring nodes. This grid can be disabled via --no-grid.

Data points provided via an optional --data-file are plotted underneath the
SOM nodes. They can be colored independent of the nodes using --data-color.

  som plot xy ... --data-file data.csv --data-color Class2`,
		Args: cobra.ExactArgs(2),
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

			_, indices, err := extractIndices(s, columns, true)
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
				tables, _, err = config.PrepareTables(reader, ignore, false, false)
				if err != nil {
					return err
				}
				data, err = extractData(config, tables, indices)
				if err != nil {
					return err
				}

				if dataColor != "" {
					colorColumn, err := reader.ReadLabels(dataColor)
					if err != nil {
						return err
					}
					dataCats, dataIndices = conv.ClassesToIndices(colorColumn, noData)
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
				classLayer := s.Layers()[indices[2][0]]
				if !classLayer.IsCategorical() {
					return fmt.Errorf("class layer %s is not categorical, can't use it for color", classLayer.Name())
				}
				classes, classIndices = conv.LayerToClasses(classLayer)
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

	command.Flags().StringVarP(&xColumn, "x-column", "x", "", "Column for x axis")
	command.Flags().StringVarP(&yColumn, "y-column", "y", "", "Column for y axis")
	command.Flags().StringVarP(&color, "color", "c", "", "Column for color")

	command.Flags().BoolVarP(&noGrid, "no-grid", "G", false, "Don't draw SOM grid lines")
	command.Flags().StringVarP(&dataColor, "data-color", "C", "", "Column for data color")

	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of the plot in pixels")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required for --labels")
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false
	command.MarkFlagFilename("data-file", "csv")

	command.MarkFlagRequired("x-column")
	command.MarkFlagRequired("y-column")

	return command
}

func extractData(conf *som.SomConfig, tables []*table.Table, indices [][2]int) (plotter.XYer, error) {
	l1, l2 := conf.Layers[indices[0][0]], conf.Layers[indices[1][0]]

	if l1.Categorical {
		return nil, fmt.Errorf("layer %s is categorical, cannot plot", l1.Name)
	}
	if l2.Categorical {
		return nil, fmt.Errorf("layer %s is categorical, cannot plot", l2.Name)
	}

	t1, t2 := tables[indices[0][0]], tables[indices[1][0]]
	c1, c2 := indices[0][1], indices[1][1]

	n1, n2 := l1.Norm[c1], l2.Norm[c2]

	xy := make([]plotter.XY, 0, tables[0].Rows())

	for i := 0; i < t1.Rows(); i++ {
		x, y := n1.DeNormalize(t1.Get(i, c1)), n2.DeNormalize(t2.Get(i, c2))
		if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
			continue
		}
		xy = append(xy, plotter.XY{
			X: x,
			Y: y,
		})
	}

	return plotter.XYs(xy), nil
}
