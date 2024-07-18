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
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func heatmapsCommand() *cobra.Command {
	var size []int
	var columns []string
	var plotColumns int

	command := &cobra.Command{
		Use:   "heatmaps [flags] <som-file> <out-file>",
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

			var layers, cols []int

			if len(columns) == 0 {
				for i, l := range s.Layers() {
					if l.IsCategorical() {
						continue
					}
					for j, c := range l.ColumnNames() {
						layers = append(layers, i)
						cols = append(cols, j)
						columns = append(columns, c)
					}
				}
			} else {
				layers, cols := make([]int, len(columns)), make([]int, len(columns))
				for i, col := range columns {
					found := false
					for j, l := range s.Layers() {
						for k, c := range l.ColumnNames() {
							if c == col {
								layers[i], cols[i] = j, k
								found = true
							}
						}
					}
					if !found {
						return fmt.Errorf("could not find column %s", col)
					}
				}
			}

			plotRows := int(math.Ceil(float64(len(columns)) / float64(plotColumns)))
			img := image.NewRGBA(image.Rect(0, 0, plotColumns*size[0], plotRows*size[1]))

			for i := range cols {
				layer, col := layers[i], cols[i]
				c, r := i%plotColumns, i/plotColumns

				subImg, err := plot.Heatmap(&s, layer, col, size[0], size[1])
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
	command.Flags().IntSliceVarP(&size, "size", "s", []int{250, 180}, "Size of the heatmap image")
	command.Flags().IntVarP(&plotColumns, "plot-columns", "p", 3, "Number of plot columns on the image")

	return command
}
