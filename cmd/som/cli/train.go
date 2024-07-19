package cli

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/yml"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

// Profiling infos:
// go tool pprof -http=":8000" -nodefraction="0.0001" som cpu.pprof

func trainCommand() *cobra.Command {
	var delim string
	var noData string
	var alpha string
	var radius string
	var epochs int
	var seed int64
	var cpuProfile bool

	command := &cobra.Command{
		Use:   "train [flags] <som-file> <data-file>",
		Short: "Trains a SOM on the given dataset",
		Long:  `Trains a SOM on the given dataset`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cpuProfile {
				stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
				defer stop.Stop()
			}

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
			trainer, err := som.NewTrainer(s, tables, trainingConfig, rand.New(rand.NewSource(seed)))
			if err != nil {
				return err
			}

			startTime := time.Now()

			progress := make(chan int)
			go func() {
				trainer.Train(epochs, progress)
			}()

			samples := tables[0].Rows()

			for epoch := range progress {
				s := samples * epoch
				samplesPerSec := float64(s) / time.Since(startTime).Seconds()
				fmt.Fprintf(os.Stderr, "\rEpoch %6d / %6d, (avg. %6d samples/sec)", epoch, epochs, int(samplesPerSec))
			}
			fmt.Fprintf(os.Stderr, "\n")

			outYaml, err := yml.ToYAML(s)
			if err != nil {
				return err
			}
			fmt.Println(string(outYaml))

			return nil
		},
	}

	command.Flags().StringVarP(&alpha, "alpha", "a", "polynomial 0.9 0.01 2",
		`Learning rate function. Options:
  - linear <start> <end>
  - power <start> <end>
  - polynomial <start> <end> <exp>
   `)
	command.Flags().StringVarP(&radius, "radius", "r", "polynomial 10 0.7 2", "Radius function. Same options as alpha")

	command.Flags().IntVarP(&epochs, "epochs", "e", 1000, "Number of epochs")
	command.Flags().Int64VarP(&seed, "seed", "s", 42, "Random seed")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "", "No data string")

	command.Flags().BoolVar(&cpuProfile, "profile", false, "Enable CPU profiling")

	command.Flags().SortFlags = false

	return command
}
