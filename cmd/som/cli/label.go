package cli

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func labelCommand() *cobra.Command {
	var delim string
	var noData string
	var column string
	var seed int64

	command := &cobra.Command{
		Use:   "label [flags] <som-file> <data-file>",
		Short: "Classifies SOM nodes using label propagation",
		Long:  `Classifies SOM nodes using label propagation`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			dataFile := args[1]

			somYaml, err := os.ReadFile(somFile)
			if err != nil {
				return err
			}
			config, trainingConfig, err := yml.ToSomConfig(somYaml)
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
			tables, err := config.PrepareTables(reader, true)
			if err != nil {
				return err
			}
			labels, err := reader.ReadLabels(column)
			if err != nil {
				return err
			}
			classes, indices := conv.ClassesToIndices(labels)

			s, err := som.New(config)
			if err != nil {
				return err
			}
			trainer, err := som.NewTrainer(s, tables, trainingConfig, rand.New(rand.NewSource(seed)))
			if err != nil {
				return err
			}

			err = trainer.PropagateLabels(column, classes, indices)
			if err != nil {
				return err
			}

			outYaml, err := yml.ToYAML(s)
			if err != nil {
				return err
			}
			fmt.Println(string(outYaml))

			return nil
		},
	}

	command.Flags().StringVarP(&column, "column", "c", "", "Column to use for label propagation")
	command.Flags().Int64VarP(&seed, "seed", "s", 42, "Random seed")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "", "No-data string")

	command.Flags().SortFlags = false
	command.MarkFlagRequired("column")

	return command
}
