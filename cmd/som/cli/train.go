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

const progressBarWidth = 36
const empty = '░'
const full = '█'

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

			tracker := newProgressTracker(epochs, tables[0].Rows())

			progress := make(chan float64, 100)
			go func() {
				trainer.Train(epochs, progress)
			}()

			epoch := 0
			for meanDist := range progress {
				tracker.Update(epoch, meanDist)
				epoch++
			}

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
	command.Flags().StringVarP(&radius, "radius", "r", "polynomial 10 0.5 2", "Radius function. Same options as alpha")

	command.Flags().IntVarP(&epochs, "epochs", "e", 1000, "Number of epochs")
	command.Flags().Int64VarP(&seed, "seed", "s", 42, "Random seed")

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "", "No data string")

	command.Flags().BoolVar(&cpuProfile, "profile", false, "Enable CPU profiling")

	command.Flags().SortFlags = false

	return command
}

type progressTracker struct {
	start   time.Time
	update  time.Time
	epochs  int
	samples int
	bar     []rune
}

func newProgressTracker(epochs int, samples int) *progressTracker {
	return &progressTracker{
		start:   time.Now(),
		update:  time.Now(),
		epochs:  epochs,
		samples: samples,
		bar:     make([]rune, progressBarWidth),
	}
}

func (t *progressTracker) Update(epoch int, meanDist float64) {
	if time.Since(t.update) < 100*time.Millisecond && epoch < t.epochs-1 {
		return
	}

	s := t.samples * epoch
	barWidth := (epoch * progressBarWidth) / t.epochs

	for i := range t.bar {
		if i < barWidth {
			t.bar[i] = full
		} else {
			t.bar[i] = empty
		}
	}
	samplesPerSec := float64(s) / time.Since(t.start).Seconds()
	fmt.Fprintf(os.Stderr, "\r[%s] %6d samples/sec | δ %5.2f", string(t.bar), int(samplesPerSec), meanDist)

	t.update = time.Now()
}

func (t *progressTracker) Finish() {
	fmt.Fprintf(os.Stderr, "\n")
}
