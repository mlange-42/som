package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/table"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func fillCommand() *cobra.Command {
	var delim string
	var noData string
	var preserve []string
	var ignore []string

	command := &cobra.Command{
		Use:   "fill [flags] <som-file> <data-file>",
		Short: "Fills missing data in the data file based on a trained SOM.",
		Long:  `Fills missing data in the data file based on a trained SOM.`,
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

			tables, original, err := config.PrepareTables(reader, ignore, false, true)
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

			err = pred.FillMissing(original)
			if err != nil {
				return err
			}

			return writeResultTables(s, original, preserve, preserved, del[0], noData)
		},
	}
	command.Flags().StringSliceVarP(&preserve, "preserve", "p", nil, "Preserve columns and prepend them to the output table")
	command.Flags().StringSliceVarP(&ignore, "ignore", "i", []string{}, "Ignore these layers for BMU search")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data string")

	command.Flags().SortFlags = false

	return command
}

func writeResultTables(s *som.Som, tables []*table.Table, preserve []string, preserved [][]string, delim rune, noData string) error {
	tabs := []*table.Table{}

	for i, t := range tables {
		if t == nil {
			continue
		}

		lay := s.Layers()[i]
		if !lay.IsCategorical() {
			tabs = append(tabs, t)
			continue
		}

		names, indices := conv.TableToClasses(t)
		lab := make([]string, len(indices))
		for j, idx := range indices {
			lab[j] = names[idx]
		}
		preserved = append(preserved, lab)
		preserve = append(preserve, lay.Name())
	}

	writer := strings.Builder{}
	err := csv.TablesToCsv(tabs, preserve, preserved, &writer, delim, noData)

	if err != nil {
		return err
	}

	fmt.Println(writer.String())
	return nil
}
