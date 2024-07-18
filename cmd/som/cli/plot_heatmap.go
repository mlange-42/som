package cli

import (
	"fmt"
	"image/png"
	"os"
	"path"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func heatmapCommand() *cobra.Command {
	var size []int
	var column string

	command := &cobra.Command{
		Use:   "heatmap [flags] <som-file> <out-file>",
		Short: "Plots heat maps of SOM variables",
		Long:  `Plots heat maps of SOM variables`,
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

			layer, col := -1, -1
			for i, l := range s.Layers() {
				for j, c := range l.ColumnNames() {
					if c == column {
						layer, col = i, j
					}
				}
			}
			if layer == -1 || col == -1 {
				return fmt.Errorf("could not find column %s", column)
			}

			img, err := plot.Heatmap(&s, layer, col, size[0], size[1], nil, nil)
			if err != nil {
				return err
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

	command.Flags().StringVarP(&column, "column", "c", "", "Column to plot")
	command.Flags().IntSliceVarP(&size, "size", "s", []int{250, 180}, "Size of the heatmap image")

	return command
}
