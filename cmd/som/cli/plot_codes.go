package cli

import (
	"context"
	"fmt"
	"image"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
	"github.com/spf13/cobra"
)

func plotCodesCommand() *cobra.Command {
	cliArgs := codePlotArgs{}

	command := &cobra.Command{
		Use:   "codes [flags] <som-file> <out-file>",
		Short: "Plots SOM nodes codes in different ways",
		Long:  `Plots SOM nodes codes in different ways`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			cliArgs.OutFile = args[1]

			if len(cliArgs.Size) != 2 {
				return fmt.Errorf("size must be two integers")
			}

			_, s, err := readSom(somFile)
			if err != nil {
				return err
			}

			cliArgs.Som = s

			ctx := context.WithValue(cmd.Context(), codePlotKey{}, cliArgs)
			cmd.SetContext(ctx)

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	command.PersistentFlags().StringSliceVarP(&cliArgs.Columns, "columns", "c", nil, "Columns to use for the heatmap (default all)")
	command.PersistentFlags().BoolVarP(&cliArgs.Normalized, "normalized", "n", false, "Use raw, normalized node weights")
	command.PersistentFlags().IntSliceVarP(&cliArgs.Size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")
	command.PersistentFlags().SortFlags = false

	command.AddCommand(plotCodesLinesCommand())

	return command
}

type codePlotArgs struct {
	Som        *som.Som
	OutFile    string
	Columns    []string
	Normalized bool
	Size       []int
}

type codePlotKey struct{}

func plotCodesLinesCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "line [flags] <som-file> <out-file>",
		Short: "Plots SOM nodes codes as line charts",
		Long:  `Plots SOM nodes codes as line charts`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs, ok := cmd.Context().Value(codePlotKey{}).(codePlotArgs)
			if !ok {
				return fmt.Errorf("args not found in context")
			}

			_, indices, err := extractIndices(cliArgs.Som, cliArgs.Columns, false)
			if err != nil {
				return err
			}

			plotType := plot.CodeLines{}
			img, err := plot.Codes(cliArgs.Som, indices, cliArgs.Normalized, &plotType, image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
			if err != nil {
				return err
			}

			return writeImage(img, cliArgs.OutFile)
		},
	}
	return command
}
