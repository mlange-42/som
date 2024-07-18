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
		Use:   "export",
		Short: "Exports a SOM to a CSV table of node vectors",
		Long:  `Exports a SOM to a CSV table of node vectors`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]

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

			s, err := som.New(config)
			if err != nil {
				return err
			}

			writer := strings.Builder{}
			err = csv.SomToCsv(&s, &writer, del[0], noData)
			if err != nil {
				return err
			}

			fmt.Println(writer.String())

			return nil
		},
	}

	command.Flags().StringVarP(&delim, "delimiter", "d", ",", "CSV delimiter")
	command.Flags().StringVarP(&noData, "no-data", "n", "-", "No data string")

	return command
}
