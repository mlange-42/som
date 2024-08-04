package cli

import (
	"fmt"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/table"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

func fillCommand() *cobra.Command {
	var delim string
	var noData string
	var preserve []string
	var ignore []string
	var kdTree bool

	var cpuProfile bool

	command := &cobra.Command{
		Use:   "fill [flags] <som-file> <data-file>",
		Short: "Fills missing data in the data file based on a trained SOM.",
		Long: `Fills missing data in the data file based on a trained SOM.

A table with the same row order as the input table is created.
Values matching the string provided via --no-data are filled based on the SOM.
Columns in the output table are the those of the input table
that also appear as a variable in the SOM.

Further columns from the input table can be transferred to the output table by use
of the --preserve flag. Here is how to transfer 'ID' and 'Name' columns:

  sum fill som.yml data.csv --preserve ID,Name

The result table is written to STDOUT in CSV format.
Redirect output to a file like this:
 
  som fill som.yml data.csv > filled.csv`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cpuProfile {
				stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
				defer stop.Stop()
			}

			somFile := args[0]
			dataFile := args[1]

			config, _, err := readConfig(somFile, false)
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
			pred, err := som.NewPredictor(s, tables, kdTree)
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
	command.Flags().BoolVarP(&kdTree, "kd-tree", "k", false, "Use kd-tree accelerated BMU search")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No-data string")

	command.Flags().BoolVar(&cpuProfile, "profile", false, "Enable CPU profiling")

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
