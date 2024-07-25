package cli

import (
	"context"
	"fmt"
	"image"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

var stepStyles = map[string]plotter.StepKind{
	"none": plotter.NoStep,
	"mid":  plotter.MidStep,
	"pre":  plotter.PreStep,
	"post": plotter.PostStep,
}

type codePlotArgs struct {
	Som        *som.Som
	OutFile    string
	Columns    []string
	Boundaries string
	Normalized bool
	ZeroAxis   bool
	Size       []int
}

type codePlotKey struct{}

func plotCodesCommand() *cobra.Command {
	cliArgs := codePlotArgs{}

	command := &cobra.Command{
		Use:   "codes [command]",
		Short: "Plots SOM node codes in different ways. See sub-commands.",
		Long: `Plots SOM node codes in different ways. See sub-commands.

All sub-commands create an image that shows all SOM nodes,
with a small plot representing each node. Each sub-command uses
a different type of plot for the nodes.

SOM variables to show in each plot can be restricted using --columns
By default, all non-categorical variables are used.

For SOMs with categorical variables, --boundaries can be used to show
boundaries between categories.`,
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
	command.PersistentFlags().StringVarP(&cliArgs.Boundaries, "boundaries", "b", "", "Optional categorical variable to show boundaries for")
	command.PersistentFlags().BoolVarP(&cliArgs.Normalized, "normalized", "n", false, "Use raw, normalized node weights")
	command.PersistentFlags().BoolVarP(&cliArgs.ZeroAxis, "zero", "z", false, "Zero the y-axis lower limit")

	command.PersistentFlags().IntSliceVarP(&cliArgs.Size, "size", "s", []int{600, 400}, "Size of the plot in pixels")
	command.PersistentFlags().SortFlags = false

	command.AddCommand(plotCodesLinesCommand())
	command.AddCommand(plotCodesBarsCommand())
	command.AddCommand(plotCodesPiesCommand())
	command.AddCommand(plotCodesRoseCommand())
	command.AddCommand(plotCodesImageCommand())

	addTreeToHelp(command, false)

	return command
}

func plotCodesLinesCommand() *cobra.Command {
	var stepStyle string
	var vertical bool
	var autoAxis bool

	command := &cobra.Command{
		Use:   "line [flags] <som-file> <out-file>",
		Short: "Plots SOM node codes as line charts.",
		Long: `Plots SOM node codes as line charts.

Create an image that shows all SOM nodes, with a small line chart
representing each node. Useful for visualizing data that is usually
shown as line charts, like time series or similarly ordered data.

Vertical line charts can be created by setting --vertical.
The step style (none, mid, pre, post) of the line charts can be set with --step.

By default, the y-axis is automatically adjusted to fit the data range
of all nodes. This behavior can be disabled with --auto,
so that each plot uses its individual axis range.
(Applies to the x-axis for vertical plots)

SOM variables to show in each plot can be restricted using --columns
By default, all non-categorical variables are used.

For SOMs with categorical variables, --boundaries can be used to show
boundaries between categories.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs, ok := cmd.Context().Value(codePlotKey{}).(codePlotArgs)
			if !ok {
				return fmt.Errorf("args not found in context")
			}

			_, indices, err := extractIndices(cliArgs.Som, cliArgs.Columns, true, false)
			if err != nil {
				return err
			}

			boundIndex := -1
			if cliArgs.Boundaries != "" {
				var err error
				_, boundIndices, err := extractIndices(cliArgs.Som, []string{cliArgs.Boundaries}, false, true)
				if err != nil {
					return err
				}
				boundIndex = boundIndices[0][0]
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
			img, err := plot.Codes(cliArgs.Som, indices, boundIndex,
				cliArgs.Normalized, cliArgs.ZeroAxis, &plotType,
				image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
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

func plotCodesBarsCommand() *cobra.Command {
	var colors []string
	var vertical bool
	var autoAxis bool

	command := &cobra.Command{
		Use:   "bar [flags] <som-file> <out-file>",
		Short: "Plots SOM node codes as bar charts.",
		Long: `Plots SOM node codes as bar charts.

Create an image that shows all SOM nodes, with a small bar chart
representing each node. Useful for visualizing data that is usually
shown as bar charts, like proportional or time series data.

Colors of the bars can be customized using --colors.

Vertical bar charts can be created by setting --vertical.
The step style (none, mid, pre, post) of the line charts can be set with --step.

By default, the y-axis is automatically adjusted to fit the data range
of all nodes. This behavior can be disabled with --auto,
so that each plot uses its individual axis range.
(Applies to the x-axis for vertical plots)

SOM variables to show in each plot can be restricted using --columns
By default, all non-categorical variables are used.

For SOMs with categorical variables, --boundaries can be used to show
boundaries between categories.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs, ok := cmd.Context().Value(codePlotKey{}).(codePlotArgs)
			if !ok {
				return fmt.Errorf("args not found in context")
			}

			_, indices, err := extractIndices(cliArgs.Som, cliArgs.Columns, true, false)
			if err != nil {
				return err
			}

			boundIndex := -1
			if cliArgs.Boundaries != "" {
				var err error
				_, boundIndices, err := extractIndices(cliArgs.Som, []string{cliArgs.Boundaries}, false, true)
				if err != nil {
					return err
				}
				boundIndex = boundIndices[0][0]
			}

			cols, err := stringsToColors(colors)
			if err != nil {
				return err
			}

			plotType := plot.CodeBar{
				Colors:     cols,
				Horizontal: vertical,
				AdjustAxis: !autoAxis,
			}

			img, err := plot.Codes(cliArgs.Som, indices, boundIndex,
				cliArgs.Normalized, cliArgs.ZeroAxis, &plotType,
				image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
			if err != nil {
				return err
			}

			return writeImage(img, cliArgs.OutFile)
		},
	}

	command.Flags().StringSliceVarP(&colors, "colors", "C", nil, "Colors for pie slices")
	command.Flags().BoolVarP(&vertical, "vertical", "v", false, "Plot bars arranged vertically (i.e. horizontal bars)")
	command.PersistentFlags().BoolVarP(&autoAxis, "auto", "a", false, "Automatically scale sub-plot axes, individually")

	return command
}

func plotCodesPiesCommand() *cobra.Command {
	var colors []string

	command := &cobra.Command{
		Use:   "pie [flags] <som-file> <out-file>",
		Short: "Plots SOM node codes as pie charts.",
		Long: `Plots SOM node codes as pie charts.

Create an image that shows all SOM nodes, with a small pie chart
representing each node. Useful for proportional data like percentages.

Colors of the pie slices can be customized using --colors.

SOM variables to show in each plot can be restricted using --columns
By default, all non-categorical variables are used.

For SOMs with categorical variables, --boundaries can be used to show
boundaries between categories.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs, ok := cmd.Context().Value(codePlotKey{}).(codePlotArgs)
			if !ok {
				return fmt.Errorf("args not found in context")
			}

			_, indices, err := extractIndices(cliArgs.Som, cliArgs.Columns, true, false)
			if err != nil {
				return err
			}

			boundIndex := -1
			if cliArgs.Boundaries != "" {
				var err error
				_, boundIndices, err := extractIndices(cliArgs.Som, []string{cliArgs.Boundaries}, false, true)
				if err != nil {
					return err
				}
				boundIndex = boundIndices[0][0]
			}

			cols, err := stringsToColors(colors)
			if err != nil {
				return err
			}

			plotType := plot.CodePie{
				Colors: cols,
			}
			img, err := plot.Codes(cliArgs.Som, indices, boundIndex,
				cliArgs.Normalized, cliArgs.ZeroAxis, &plotType,
				image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
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
		Long: `Plots SOM node codes as rose alias Nightingale charts.

Create an image that shows all SOM nodes, with a small rose or Nightingale chart
representing each node. Useful for many different kinds of data.

Colors of the rose slices can be customized using --colors.

SOM variables to show in each plot can be restricted using --columns
By default, all non-categorical variables are used.

For SOMs with categorical variables, --boundaries can be used to show
boundaries between categories.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs, ok := cmd.Context().Value(codePlotKey{}).(codePlotArgs)
			if !ok {
				return fmt.Errorf("args not found in context")
			}

			_, indices, err := extractIndices(cliArgs.Som, cliArgs.Columns, true, false)
			if err != nil {
				return err
			}
			boundIndex := -1
			if cliArgs.Boundaries != "" {
				var err error
				_, boundIndices, err := extractIndices(cliArgs.Som, []string{cliArgs.Boundaries}, false, true)
				if err != nil {
					return err
				}
				boundIndex = boundIndices[0][0]
			}

			cols, err := stringsToColors(colors)
			if err != nil {
				return err
			}

			plotType := plot.CodeRose{
				Colors: cols,
			}
			img, err := plot.Codes(cliArgs.Som, indices, boundIndex,
				cliArgs.Normalized, cliArgs.ZeroAxis, &plotType,
				image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
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
		Long: `Plots SOM node codes as images.

Create an image that shows all SOM nodes, with a small image or heatmap
representing each node. Useful for image-like data and other gridded data.

The size of the image id controlled using --rows. The number of image columns
is determined automatically based on the number of variables and rows.

SOM variables to show in each plot can be restricted using --columns
By default, all non-categorical variables are used.

For SOMs with categorical variables, --boundaries can be used to show
boundaries between categories.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs, ok := cmd.Context().Value(codePlotKey{}).(codePlotArgs)
			if !ok {
				return fmt.Errorf("args not found in context")
			}

			_, indices, err := extractIndices(cliArgs.Som, cliArgs.Columns, true, false)
			if err != nil {
				return err
			}

			boundIndex := -1
			if cliArgs.Boundaries != "" {
				var err error
				_, boundIndices, err := extractIndices(cliArgs.Som, []string{cliArgs.Boundaries}, false, true)
				if err != nil {
					return err
				}
				boundIndex = boundIndices[0][0]
			}

			plotType := plot.CodeImage{
				Rows: rows,
			}
			img, err := plot.Codes(cliArgs.Som, indices, boundIndex,
				cliArgs.Normalized, cliArgs.ZeroAxis, &plotType,
				image.Pt(cliArgs.Size[0], cliArgs.Size[1]))
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
