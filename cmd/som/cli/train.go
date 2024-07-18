package cli

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func trainCommand() *cobra.Command {
	var delim string
	var noData string
	var alpha string
	var radius string
	var epochs int
	var seed int64

	command := &cobra.Command{
		Use:   "train [flags] <som-file> <data-file>",
		Short: "Trains a SOM on the given dataset",
		Long:  `Trains a SOM on the given dataset`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			dataFile := args[1]

			somYaml, err := os.ReadFile(somFile)
			if err != nil {
				return err
			}
			config, err := yml.ToSomConfig(somYaml)
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

			learningDecay, err := decay.FromString(alpha)
			if err != nil {
				return err
			}
			radiusDecay, err := decay.FromString(radius)
			if err != nil {
				return err
			}

			trainingConfig := &som.TrainingConfig{
				LearningRate:       learningDecay,
				NeighborhoodRadius: radiusDecay,
			}

			s, err := som.New(config)
			if err != nil {
				return err
			}
			trainer, err := som.NewTrainer(&s, tables, trainingConfig, rand.New(rand.NewSource(seed)))
			if err != nil {
				return err
			}

			trainer.Train(epochs)

			outYaml, err := yml.ToYAML(&s)
			if err != nil {
				return err
			}
			fmt.Println(string(outYaml))

			return nil
		},
	}

	command.Flags().StringVarP(&alpha, "alpha", "a", "linear 0.5 0.01", "Learning rate function")
	command.Flags().StringVarP(&radius, "radius", "r", "linear 10 0.5", "Radius function")

	command.Flags().IntVarP(&epochs, "epochs", "e", 1000, "Number of epochs")
	command.Flags().Int64VarP(&seed, "seed", "s", 42, "Random seed")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "-", "No data string")

	return command
}
