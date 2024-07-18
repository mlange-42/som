package cli

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func heatmapCommand() *cobra.Command {
	var size []int
	var columns []string
	var plotColumns int
	var dataFile string
	var labelsColumn string
	var delim string
	var noData string

	command := &cobra.Command{
		Use:   "heatmap [flags] <som-file> <out-file>",
		Short: "Plots heat maps of multiple SOM variables",
		Long:  `Plots heat maps of multiple SOM variables`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			outFile := args[1]

			somYaml, err := os.ReadFile(somFile)
			if err != nil {
				return err
			}
			config, err := yml.ToSomConfig(somYaml)
			if err != nil {
				return err
			}

			s, err := som.New(config)
			if err != nil {
				return err
			}

			columns, indices, err := extractIndices(&s, columns)
			if err != nil {
				return err
			}

			if plotColumns == 0 {
				plotColumns = int(math.Sqrt(float64(len(columns))))
			}

			plotRows := int(math.Ceil(float64(len(columns)) / float64(plotColumns)))
			img := image.NewRGBA(image.Rect(0, 0, plotColumns*size[0], plotRows*size[1]))

			var labels []string
			var positions []plotter.XY

			if labelsColumn != "" {
				if dataFile == "" {
					return fmt.Errorf("data file must be specified when labels column is specified")
				}

				del := []rune(delim)
				if len(delim) != 1 {
					panic("delimiter must be a single character")
				}
				labels, positions, err = extractLabels(config, &s, labelsColumn, dataFile, del[0], noData)
				if err != nil {
					return err
				}
			}

			for i := range indices {
				layer, col := indices[i][0], indices[i][1]
				c, r := i%plotColumns, i/plotColumns

				l := &s.Layers()[layer]
				grid := plot.SomLayerGrid{Som: &s, Layer: layer, Column: col}
				title := fmt.Sprintf("%s: %s", l.Name(), l.ColumnNames()[col])

				subImg, err := plot.Heatmap(title, &grid, size[0], size[1], labels, positions)
				if err != nil {
					return err
				}

				draw.Draw(img, image.Rect(c*size[0], r*size[1], (c+1)*size[0], (r+1)*size[1]), subImg, image.Point{}, draw.Src)
			}

			err = os.MkdirAll(path.Dir(outFile), os.ModePerm)
			if err != nil {
				return err
			}
			file, err := os.Create(outFile)
			if err != nil {
				return err
			}

			return png.Encode(file, img)
		},
	}

	command.Flags().StringSliceVarP(&columns, "columns", "c", nil, "Column names for the heatmap")
	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")
	command.Flags().IntVarP(&plotColumns, "plot-columns", "p", 0, "Number of plot columns on the image")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file")
	command.Flags().StringVarP(&labelsColumn, "labels", "l", "", "Labels column in the data file")
	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "-", "No.data value")

	return command
}

func extractIndices(s *som.Som, columns []string) ([]string, [][2]int, error) {
	var indices [][2]int

	if len(columns) == 0 {
		for i, l := range s.Layers() {
			if l.IsCategorical() {
				continue
			}
			for j, c := range l.ColumnNames() {
				indices = append(indices, [2]int{i, j})
				columns = append(columns, c)
			}
		}
	} else {
		indices = make([][2]int, len(columns))
		for i, col := range columns {
			found := false
			for j, l := range s.Layers() {
				for k, c := range l.ColumnNames() {
					if c == col {
						indices[i] = [2]int{j, k}
						found = true
					}
				}
			}
			if !found {
				return nil, nil, fmt.Errorf("could not find column %s", col)
			}
		}
	}

	return columns, indices, nil
}

func extractLabels(
	config *som.SomConfig, s *som.Som,
	labelsColumn string, dataFile string, delim rune, noData string) ([]string, []plotter.XY, error) {

	reader, err := csv.NewFileReader(dataFile, delim, noData)
	if err != nil {
		return nil, nil, err
	}

	labels, err := reader.ReadLabels(labelsColumn)
	if err != nil {
		return nil, nil, err
	}

	tables, err := config.PrepareTables(reader, false)
	if err != nil {
		return nil, nil, err
	}

	pred, err := som.NewPredictor(s, tables)
	if err != nil {
		return nil, nil, err
	}

	bmu, err := pred.GetBMU()
	if err != nil {
		return nil, nil, err
	}

	perCell := make([]int, s.Size().Height*s.Size().Width)
	for i := range labels {
		idx := int(bmu.Get(i, 0))
		perCell[idx]++
	}
	count := make([]int, s.Size().Height*s.Size().Width)

	xy := make([]plotter.XY, len(labels))
	for i := range labels {
		idx := int(bmu.Get(i, 0))

		frac := float64(count[idx]+1) / float64(perCell[idx]+1)

		xy[i].X, xy[i].Y = bmu.Get(i, 1), bmu.Get(i, 2)-0.5+frac
		count[idx]++
	}
	return labels, xy, nil
}
