package cli

import (
	"fmt"
	"image"

	"github.com/mlange-42/som/plot"
	"github.com/spf13/cobra"
)

func plotCodesCommand() *cobra.Command {
	var size []int
	var columns []string

	command := &cobra.Command{
		Use:   "codes [flags] <som-file> <out-file>",
		Short: "Plots SOM nodes codes in different ways",
		Long:  `Plots SOM nodes codes in different ways`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			outFile := args[1]

			if len(size) != 2 {
				return fmt.Errorf("size must be two integers")
			}

			_, s, err := readSom(somFile)
			if err != nil {
				return err
			}

			_, indices, err := extractIndices(s, columns, false)
			if err != nil {
				return err
			}

			img := plot.Codes(s, indices, image.Pt(size[0], size[1]))

			return writeImage(img, outFile)
		},
	}

	command.Flags().StringSliceVarP(&columns, "columns", "c", nil, "Columns to use for the heatmap (default all)")
	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")

	command.Flags().SortFlags = false

	return command
}
