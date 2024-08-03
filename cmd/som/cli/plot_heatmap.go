package cli

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/table"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/plotter"
)

func plotHeatmapCommand() *cobra.Command {
	var size []int
	var columns []string
	var boundaries string
	var plotColumns int
	var dataFile string
	var labelsColumn string
	var delim string
	var noData string
	var ignore []string
	var sample int

	command := &cobra.Command{
		Use:   "heatmap [flags] <som-file> <out-file>",
		Short: "Plots heat maps of multiple SOM variables, a.k.a. components plot.",
		Long: `Plots heat maps of multiple SOM variables, a.k.a. components plot.

By default, the command creates an image with multiple panels.
Each panel shows a heatmap for one of the SOM's variables.
Categorical variables are converted to "class maps" with a unique
color for each category.

To select individual variables or a sub-set of variables, use --columns.

For SOMs with categorical variables, --boundaries can be used to show
boundaries between categories.

Data provided via --data-file can be displayed on top of the heatmaps,
showing the values in the column given by the --label flag:

  som plot heatmap som.yml heatmap.png --data-file data.csv --label name

For large datasets, --sample can be used to show only a sub-set of the data.`,
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

			columns, indices, err := extractIndices(s, columns, true, true)
			if err != nil {
				return err
			}

			if plotColumns == 0 {
				plotColumns = int(math.Sqrt(float64(len(columns))))
			}

			plotRows := int(math.Ceil(float64(len(columns)) / float64(plotColumns)))
			img := image.NewRGBA(image.Rect(0, 0, plotColumns*size[0], plotRows*size[1]))

			var reader table.Reader
			var predictor *som.Predictor

			if dataFile != "" {
				var err error
				reader, err = csv.NewFileReader(dataFile, del[0], noData)
				if err != nil {
					return err
				}
				predictor, _, err = createPredictor(config, s, reader, ignore)
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

				labels, positions, err = extractLabels(predictor, labelsColumn, reader, sample)
				if err != nil {
					return err
				}
			}

			bounds, err := extractBoundariesLayer(s, boundaries)
			if err != nil {
				return err
			}

			for i := range indices {
				layer, col := indices[i][0], indices[i][1]
				c, r := i%plotColumns, i/plotColumns

				title, classes, grid := createTitleAndGrid(s, layer, col)
				subImg, err := plot.Heatmap(title, grid, bounds, size[0], size[1], classes, labels, positions)
				if err != nil {
					return err
				}

				draw.Draw(img, image.Rect(c*size[0], r*size[1], (c+1)*size[0], (r+1)*size[1]), subImg, image.Point{}, draw.Src)
			}

			return writeImage(img, outFile)
		},
	}

	command.Flags().StringSliceVarP(&columns, "columns", "c", nil, "Columns to use for the heatmap (default all)")
	command.Flags().StringVarP(&boundaries, "boundaries", "b", "", "Optional categorical variable to show boundaries for")
	command.Flags().IntSliceVarP(&size, "size", "s", []int{600, 400}, "Size of individual heatmap panels")
	command.Flags().IntVarP(&plotColumns, "plot-columns", "p", 0, "Number of plot columns on the image (default sqrt(#cols))")
	command.Flags().StringVarP(&dataFile, "data-file", "f", "", "Data file. Required for --label")
	command.Flags().StringVarP(&labelsColumn, "label", "l", "", "Label column in the data file")
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")
	command.Flags().IntVarP(&sample, "sample", "S", 0, "Sample this many rows from the data file (default all)")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data value (default \"\")")

	command.Flags().SortFlags = false

	command.MarkFlagFilename("data-file", "csv")

	return command
}

func createTitleAndGrid(s *som.Som, layer, col int) (title string, classes []string, grid plotter.GridXYZ) {
	l := s.Layers()[layer]
	if col >= 0 {
		grid = &plot.SomLayerGrid{Som: s, Layer: layer, Column: col}
		title = fmt.Sprintf("%s: %s", l.Name(), l.ColumnNames()[col])
	} else {
		title = l.Name()
		var classIndices []int
		classes, classIndices = conv.LayerToClasses(l)
		grid = &plot.IntGrid{Size: *s.Size(), Values: classIndices}
	}
	return
}

func extractBoundariesLayer(s *som.Som, boundaries string) (plotter.GridXYZ, error) {
	var bounds plotter.GridXYZ
	if boundaries != "" {
		_, idx, err := extractIndices(s, []string{boundaries}, false, true)
		if err != nil {
			return nil, err
		}
		_, classIndices := conv.LayerToClasses(s.Layers()[idx[0][0]])
		bounds = &plot.IntGrid{Size: *s.Size(), Values: classIndices}
	}
	return bounds, nil
}

func writeImage(img image.Image, outFile string) error {
	err := os.MkdirAll(path.Dir(outFile), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func readSom(somFile string) (*som.SomConfig, *som.Som, error) {
	somYaml, err := os.ReadFile(somFile)
	if err != nil {
		return nil, nil, err
	}
	config, _, err := yml.ToSomConfig(somYaml)
	if err != nil {
		return nil, nil, err
	}

	s, err := som.New(config)
	if err != nil {
		return nil, nil, err
	}

	return config, s, nil
}

func extractIndices(s *som.Som, columns []string, inclContinuous, inclCategorical bool) ([]string, [][2]int, error) {
	var indices [][2]int

	if len(columns) == 0 {
		return extractAllIndices(s, inclContinuous, inclCategorical)
	}

	indices = make([][2]int, len(columns))
	for i, col := range columns {
		found := false
		for j, l := range s.Layers() {
			if l.IsCategorical() {
				if col == l.Name() {
					if !inclCategorical {
						return nil, nil, fmt.Errorf("column %s is in categorical layer %s but categorical layers are excluded", col, l.Name())
					}
					indices[i] = [2]int{j, -1}
					found = true
					break
				}
				continue
			}
			for k, c := range l.ColumnNames() {
				if c == col {
					if !inclContinuous {
						return nil, nil, fmt.Errorf("column %s is in continuous layer %s but continuous layers are excluded", col, l.Name())
					}
					indices[i] = [2]int{j, k}
					found = true
					break
				}
			}
		}
		if !found {
			return nil, nil, fmt.Errorf("could not find column %s", col)
		}
	}

	return columns, indices, nil
}

func extractAllIndices(s *som.Som, inclContinuous, inclCategorical bool) (columns []string, indices [][2]int, err error) {
	for i, l := range s.Layers() {
		if l.IsCategorical() {
			if inclCategorical {
				indices = append(indices, [2]int{i, -1})
				columns = append(columns, l.Name())
			}
			continue
		}
		if inclContinuous {
			for j, c := range l.ColumnNames() {
				indices = append(indices, [2]int{i, j})
				columns = append(columns, c)
			}
		}
	}
	return columns, indices, nil
}

func createPredictor(config *som.SomConfig, s *som.Som, reader table.Reader, ignoreLayers []string) (*som.Predictor, []*table.Table, error) {
	tables, _, err := config.PrepareTables(reader, ignoreLayers, false, false)
	if err != nil {
		return nil, nil, err
	}

	pred, err := som.NewPredictor(s, tables, false)
	if err != nil {
		return nil, nil, err
	}

	return pred, tables, nil
}

func extractLabels(predictor *som.Predictor,
	labelsColumn string, reader table.Reader, sample int) ([]string, []plotter.XY, error) {

	labels, err := reader.ReadLabels(labelsColumn)
	if err != nil {
		return nil, nil, err
	}

	bmu := predictor.GetBMUTable()
	nodes := predictor.Som().Size().Nodes()

	indices := make([]int, len(labels))
	for i := range indices {
		indices[i] = i
	}

	var outLabel []string
	if sample > 0 && sample < len(indices) {
		rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })
		indices = indices[:sample]
		outLabel = make([]string, sample)
	}

	perCell := make([]int, nodes)
	for _, i := range indices {
		idx := int(bmu.Get(i, 0))
		perCell[idx]++
	}
	count := make([]int, nodes)

	xy := make([]plotter.XY, len(indices))
	for c, i := range indices {
		idx := int(bmu.Get(i, 0))

		frac := float64(count[idx]+1) / float64(perCell[idx]+1)

		xy[c].X, xy[c].Y = bmu.Get(i, 1), bmu.Get(i, 2)-0.5+frac
		count[idx]++

		if outLabel != nil {
			outLabel[c] = labels[i]
		}
	}
	if outLabel != nil {
		return outLabel, xy, nil
	}
	return labels, xy, nil
}
