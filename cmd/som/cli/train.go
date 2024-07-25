package cli

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/decay"
	"github.com/mlange-42/som/table"
	"github.com/mlange-42/som/yml"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	var decayFunc string
	var epochs int
	var seed int64
	var cpuProfile bool
	var visomLambda float64
	var progressFile string
	var progressInterval int

	var command *cobra.Command
	command = &cobra.Command{
		Use:   "train [flags] <som-file> <data-file>",
		Short: "Trains an SOM on the given dataset.",
		Long: `Trains an SOM on the given dataset.

The trained SOM is written to STDOUT in YAML format.
Redirect output to a file like this:

  som train som.yml data.csv > trained.yml

Learning parameters are usually specified in the SOM's YAML file,
but can also be set or overwritten using the provided CLI flags.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cpuProfile {
				stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
				defer stop.Stop()
			}

			somFile := args[0]
			dataFile := args[1]

			config, trainingConfig, err := readConfig(somFile)
			if err != nil {
				return err
			}

			del := []rune(delim)
			if len(delim) != 1 {
				return fmt.Errorf("delimiter must be a single character")
			}

			tables, err := prepareTables(config, dataFile, del[0], noData)
			if err != nil {
				return err
			}

			err = overwriteParameters(command, trainingConfig,
				epochs, visomLambda, alpha, radius, decayFunc)
			if err != nil {
				return err
			}

			s, err := runTraining(config, trainingConfig, tables, seed, progressFile, progressInterval, del[0])
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

	command.Flags().StringVarP(&alpha, "alpha", "a", "polynomial 0.25 0.01 2",
		`Overwrites the learning rate function of the SOM file.
Options:
  - linear <start> <end>
  - power <start> <end>
  - polynomial <start> <end> <exp>
   `)
	command.Flags().StringVarP(&radius, "radius", "r", "polynomial 10 0.7 2", "Overwrites the radius function of the SOM file.\nSame options as alpha")
	command.Flags().StringVarP(&decayFunc, "decay", "d", "", "Overwrites the weight decay function of the SOM file.\nSame options as alpha (default no decay)")

	command.Flags().IntVarP(&epochs, "epochs", "e", 1000, "Overwrites the number of epochs of the SOM file")
	command.Flags().Int64VarP(&seed, "seed", "s", 42, "Random seed")

	command.Flags().Float64VarP(&visomLambda, "visom-lambda", "v", 0.0, "Overwrites ViSOM resolution. 0 = no ViSOM")

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No data string")

	command.Flags().IntVarP(&progressInterval, "progress", "P", 100, "Interval for progress output.\nIgnored if no <progress-file> is given")
	command.Flags().StringVarP(&progressFile, "progress-file", "p", "", "CSV file for training progress output")

	command.Flags().BoolVar(&cpuProfile, "profile", false, "Enable CPU profiling")

	command.Flags().SortFlags = false
	command.MarkFlagFilename("progress-file", "csv")

	return command
}

func overwriteParameters(command *cobra.Command, conf *som.TrainingConfig,
	epochs int, visomLambda float64, alpha, radius, decayFunc string) error {
	flagUsed := map[string]bool{}
	command.Flags().Visit(func(f *pflag.Flag) {
		flagUsed[f.Name] = true
	})

	if _, ok := flagUsed["epochs"]; ok {
		conf.Epochs = epochs
	}
	if _, ok := flagUsed["visom-lambda"]; ok {
		conf.ViSomLambda = visomLambda
	}

	var err error
	if _, ok := flagUsed["alpha"]; ok {
		conf.LearningRate, err = decay.FromString(alpha)
		if err != nil {
			return err
		}
	}
	if _, ok := flagUsed["radius"]; ok {
		conf.NeighborhoodRadius, err = decay.FromString(radius)
		if err != nil {
			return err
		}
	}
	if _, ok := flagUsed["decay"]; ok {
		if decayFunc == "" {
			conf.WeightDecay = nil
		} else {
			conf.WeightDecay, err = decay.FromString(decayFunc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runTraining(config *som.SomConfig, trainingConfig *som.TrainingConfig,
	tables []*table.Table, seed int64,
	progressFile string, writeInterval int, csvDelim rune,
) (*som.Som, error) {

	s, err := som.New(config)
	if err != nil {
		return nil, err
	}
	trainer, err := som.NewTrainer(s, tables, trainingConfig, rand.New(rand.NewSource(seed)))
	if err != nil {
		return nil, err
	}

	var writer io.Writer
	if progressFile == "" {
		writeInterval = 0
	} else {
		file, err := os.Create(progressFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		writer = file
	}
	tracker := newProgressTracker(trainingConfig.Epochs, tables[0].Rows(), writer, writeInterval, csvDelim)

	progress := make(chan som.TrainingProgress, 100)
	go trainer.Train(progress)

	epoch := 0
	for p := range progress {
		tracker.Update(epoch, &p)
		epoch++
	}

	return s, nil
}

func defaultTrainingConfig() *som.TrainingConfig {
	return &som.TrainingConfig{
		Epochs:             1000,
		LearningRate:       &decay.Polynomial{Start: 0.25, End: 0.01, Exp: 2},
		NeighborhoodRadius: &decay.Polynomial{Start: 10, End: 0.7, Exp: 2},
		ViSomLambda:        0,
	}
}

func prepareTables(config *som.SomConfig, path string, delim rune, noData string) ([]*table.Table, error) {
	reader, err := csv.NewFileReader(path, delim, noData)
	if err != nil {
		return nil, err
	}
	tables, _, err := config.PrepareTables(reader, nil, true, false)
	return tables, err
}

func readConfig(path string) (*som.SomConfig, *som.TrainingConfig, error) {
	somYaml, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	config, trainingConfig, err := yml.ToSomConfig(somYaml)
	if err != nil {
		return nil, nil, err
	}

	if trainingConfig == nil {
		trainingConfig = defaultTrainingConfig()
	}

	return config, trainingConfig, nil
}

type progressTracker struct {
	start         time.Time
	update        time.Time
	epochs        int
	samples       int
	writer        io.Writer
	writeInterval int
	csvDelim      rune
	bar           []rune
}

func newProgressTracker(epochs int, samples int, writer io.Writer, writeInterval int, csvDelim rune) *progressTracker {
	return &progressTracker{
		start:         time.Now(),
		update:        time.Now(),
		epochs:        epochs,
		samples:       samples,
		writer:        writer,
		writeInterval: writeInterval,
		csvDelim:      csvDelim,
		bar:           make([]rune, progressBarWidth),
	}
}

func (t *progressTracker) Update(epoch int, progress *som.TrainingProgress) {
	if t.writeInterval > 0 && epoch%t.writeInterval == 0 {
		if epoch == 0 {
			fmt.Fprintln(t.writer, progress.CsvHeader(t.csvDelim))
		}
		fmt.Fprintln(t.writer, progress.CsvRow(t.csvDelim))
	}

	if time.Since(t.update) < 100*time.Millisecond && epoch < t.epochs-1 {
		return
	}

	s := t.samples * (epoch + 1)
	barWidth := ((epoch + 1) * progressBarWidth) / t.epochs

	for i := range t.bar {
		if i < barWidth {
			t.bar[i] = full
		} else {
			t.bar[i] = empty
		}
	}
	samplesPerSec := float64(s) / time.Since(t.start).Seconds()
	fmt.Fprintf(os.Stderr, "\r[%s] %6d samples/sec | δ %5.2f/%5.2f", string(t.bar), int(samplesPerSec), progress.MeanDist, progress.Error)

	t.update = time.Now()
}

func (t *progressTracker) Finish() {
	fmt.Fprintf(os.Stderr, "\n")
}
