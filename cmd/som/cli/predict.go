package cli

import (
	"fmt"
	"os"
	"slices"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func predictCommand() *cobra.Command {
	var delim string
	var noData string
	var preserve []string
	var ignore []string
	var layers []string
	var writeAllLayers bool

	command := &cobra.Command{
		Use:   "predict [flags] <som-file> <data-file>",
		Short: "predict entire layers or table columns using a trained SOM",
		Long:  `predict entire layers or table columns using a trained SOM`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			dataFile := args[1]

			somYaml, err := os.ReadFile(somFile)
			if err != nil {
				return err
			}
			config, _, err := yml.ToSomConfig(somYaml)
			if err != nil {
				return err
			}

			del := []rune(delim)
			if len(delim) != 1 {
				return fmt.Errorf("delimiter must be a single character")
			}

			reader, err := csv.NewFileReader(dataFile, del[0], noData)
			if err != nil {
				return err
			}

			preserved := [][]string{}
			for _, column := range preserve {
				col, err := reader.ReadLabels(column)
				if err != nil {
					return err
				}
				preserved = append(preserved, col)
			}

			ignoreRead := append([]string{}, ignore...)
			ignoreRead = append(ignoreRead, layers...)
			tables, original, err := config.PrepareTables(reader, ignoreRead, false, true)
			if err != nil {
				return err
			}

			s, err := som.New(config)
			if err != nil {
				return err
			}
			pred, err := som.NewPredictor(s, tables)
			if err != nil {
				return err
			}

			err = pred.Predict(original, layers)
			if err != nil {
				return err
			}

			if !writeAllLayers {
				for i := range original {
					if slices.Contains(layers, s.Layers()[i].Name()) {
						continue
					}
					original[i] = nil
				}
			}

			return writeResultTables(s, original, preserve, preserved, del[0], noData)
		},
	}

	command.Flags().StringSliceVarP(&layers, "layers", "l", nil, "Predict these layers from all other layers")
	command.Flags().StringSliceVarP(&preserve, "preserve", "p", nil, "Preserve columns and prepend them to the output table")
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")
	command.Flags().BoolVarP(&writeAllLayers, "all", "a", false, "Write all layers instead of just predicted layers")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data string")

	command.Flags().SortFlags = false

	return command
}
