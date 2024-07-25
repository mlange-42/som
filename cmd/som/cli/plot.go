package cli

import (
	"fmt"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func plotCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "plot",
		Short: "Plots visualizations for an SOM in various ways. See sub-commands.",
		Long: `Plots visualizations for an SOM in various ways. See sub-commands.

All sub-commands take an SOM YAML file and an output PNG file
as positional arguments.

Please note that the built-in visualizations are not intended for
publication-quality output. Instead, they serve as quick tools for
inspecting training and prediction results. For high-quality visualizations,
we recommend exporting the SOM and other results to CSV files.
You can then use dedicated visualization libraries in languages such as
Python or R to create more refined and customized graphics.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	command.AddCommand(plotHeatmapCommand())
	command.AddCommand(plotCodesCommand())
	command.AddCommand(plotUMatrixCommand())
	command.AddCommand(plotXyCommand())
	command.AddCommand(plotDensityCommand())
	command.AddCommand(plotErrorCommand())

	addTreeToHelp(command, false)

	return command
}

func plotHeatmap(size []int,
	somFile, outFile, dataFile,
	labelsColumn, delim, noData string,
	title string,
	ignoreLayers []string, boundaries string, sampleData int,
	getData func(s *som.Som, p *som.Predictor, r table.Reader) (plotter.GridXYZ, []string, error)) error {

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

	var reader table.Reader
	var predictor *som.Predictor

	if dataFile != "" {
		var err error
		reader, err = csv.NewFileReader(dataFile, del[0], noData)
		if err != nil {
			return err
		}
		predictor, _, err = createPredictor(config, s, reader, ignoreLayers)
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
		labels, positions, err = extractLabels(predictor, labelsColumn, reader, sampleData)
		if err != nil {
			return err
		}
	}

	bounds, err := extractBoundariesLayer(s, boundaries)
	if err != nil {
		return err
	}

	grid, cats, err := getData(s, predictor, reader)
	if err != nil {
		return err
	}

	img, err := plot.Heatmap(title, grid, bounds, size[0], size[1], cats, labels, positions)
	if err != nil {
		return err
	}

	return writeImage(img, outFile)
}
