package cli

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
	"github.com/spf13/cobra"
	"golang.org/x/image/colornames"
	"gonum.org/v1/plot/plotter"
)

var stepStyles = map[string]plotter.StepKind{
	"none": plotter.NoStep,
	"mid":  plotter.MidStep,
	"pre":  plotter.PreStep,
	"post": plotter.PostStep,
}

func plotCodesCommand() *cobra.Command {
	cliArgs := codePlotArgs{}

	command := &cobra.Command{
		Use:   "codes [command]",
		Short: "Plots SOM node codes in different ways. See sub-commands.",
		Long:  `Plots SOM node codes in different ways. See sub-commands.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				_ = cmd.Help()
				return nil
			}
			if len(args) < 2 {
				return fmt.Errorf("requires sub-command")
			}
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

	command.PersistentFlags().StringSliceVarP(&cliArgs.Columns, "columns", "c", nil, "Columns to use for the codes plot (default all)")
	command.PersistentFlags().BoolVarP(&cliArgs.Normalized, "normalized", "n", false, "Use raw, normalized node weights")
	command.PersistentFlags().BoolVarP(&cliArgs.ZeroAxis, "zero", "z", false, "Zero the y-axis lower limit")

	command.PersistentFlags().IntSliceVarP(&cliArgs.Size, "size", "s", []int{600, 400}, "Size of the plot in pixels")
	command.PersistentFlags().SortFlags = false

	command.AddCommand(plotCodesLinesCommand())
	command.AddCommand(plotCodesPiesCommand())
	command.AddCommand(plotCodesRoseCommand())
	command.AddCommand(plotCodesImageCommand())

	addTreeToHelp(command, false)

	return command
}

type codePlotArgs struct {
	Som        *som.Som
	OutFile    string
	Columns    []string
	Normalized bool
	ZeroAxis   bool
	Size       []int
}

type codePlotKey struct{}

func plotCodesLinesCommand() *cobra.Command {
	var stepStyle string
	var vertical bool
	var autoAxis bool

	command := &cobra.Command{
		Use:   "line [flags] <som-file> <out-file>",
		Short: "Plots SOM node codes as line charts.",
		Long:  `Plots SOM node codes as line charts.`,
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

			step, ok := stepStyles[strings.ToLower(stepStyle)]
			if !ok {
				return fmt.Errorf("invalid step style: %s", stepStyle)
			}

			plotType := plot.CodeLines{
				StepStyle:  step,
				Vertical:   vertical,
				AdjustAxis: !autoAxis,
			}
			img, err := plot.Codes(cliArgs.Som, indices, cliArgs.Normalized, cliArgs.ZeroAxis, &plotType, image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
			if err != nil {
				return err
			}

			return writeImage(img, cliArgs.OutFile)
		},
	}

	command.Flags().StringVarP(&stepStyle, "step", "S", "none", "Line step style (none, mid, pre, post)")
	command.Flags().BoolVarP(&vertical, "vertical", "v", false, "Plot lines vertically")
	command.PersistentFlags().BoolVarP(&autoAxis, "auto", "a", false, "Automatically scale sub-plot axes, individually")

	return command
}

func plotCodesPiesCommand() *cobra.Command {
	var colors []string

	command := &cobra.Command{
		Use:   "pie [flags] <som-file> <out-file>",
		Short: "Plots SOM node codes as pie charts.",
		Long:  `Plots SOM node codes as pie charts.`,
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

			cols := make([]color.Color, len(colors))
			for i, c := range colors {
				cols[i], ok = colornames.Map[c]
				if !ok {
					return fmt.Errorf("color name %s unknown", c)
				}
			}

			plotType := plot.CodePie{
				Colors: cols,
			}
			img, err := plot.Codes(cliArgs.Som, indices, cliArgs.Normalized, cliArgs.ZeroAxis, &plotType, image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
			if err != nil {
				return err
			}

			return writeImage(img, cliArgs.OutFile)
		},
	}

	command.Flags().StringSliceVarP(&colors, "colors", "C", nil, "Colors for pie slices")

	return command
}

func plotCodesRoseCommand() *cobra.Command {
	var colors []string

	command := &cobra.Command{
		Use:   "rose [flags] <som-file> <out-file>",
		Short: "Plots SOM node codes as rose alias Nightingale charts.",
		Long:  `Plots SOM node codes as rose alias Nightingale charts.`,
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

			cols := make([]color.Color, len(colors))
			for i, c := range colors {
				cols[i], ok = colornames.Map[c]
				if !ok {
					return fmt.Errorf("color name %s unknown", c)
				}
			}

			plotType := plot.CodeRose{
				Colors: cols,
			}
			img, err := plot.Codes(cliArgs.Som, indices, cliArgs.Normalized, cliArgs.ZeroAxis, &plotType, image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
			if err != nil {
				return err
			}

			return writeImage(img, cliArgs.OutFile)
		},
	}

	command.Flags().StringSliceVarP(&colors, "colors", "C", nil, "Colors for pie slices")

	return command
}

func plotCodesImageCommand() *cobra.Command {
	var rows int

	command := &cobra.Command{
		Use:   "image [flags] <som-file> <out-file>",
		Short: "Plots SOM node codes as images.",
		Long:  `Plots SOM node codes as images.`,
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

			plotType := plot.CodeImage{
				Rows: rows,
			}
			img, err := plot.Codes(cliArgs.Som, indices, cliArgs.Normalized, cliArgs.ZeroAxis, &plotType, image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
			if err != nil {
				return err
			}

			return writeImage(img, cliArgs.OutFile)
		},
	}

	command.Flags().IntVarP(&rows, "rows", "r", 1, "Number of rows for image plot")

	command.MarkFlagRequired("rows")

	return command
}
