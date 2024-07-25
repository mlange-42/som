package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/csv"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func exportCommand() *cobra.Command {
	var delim string
	var noData string

	command := &cobra.Command{
		Use:   "export [flags] <som-file>",
		Short: "Exports an SOM to a CSV table of node vectors.",
		Long:  `Exports an SOM to a CSV table of node vectors.`,
		Args:  ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]

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

			s, err := som.New(config)
			if err != nil {
				return err
			}

			writer := strings.Builder{}
			err = csv.SomToCsv(s, &writer, del[0], noData)
			if err != nil {
				return err
			}

			fmt.Println(writer.String())

			return nil
		},
	}

	command.Flags().StringVarP(&delim, "delimiter", "D", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "N", "", "No data string")

	return command
}
