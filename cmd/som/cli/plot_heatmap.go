package cli

import (
	"image/png"
	"os"
	"path"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/plot"
	"github.com/mlange-42/som/yml"
	"github.com/spf13/cobra"
)

func heatmapCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "heatmap [flags] <som-file> <out-file>",
		Short: "Plots heat maps of SOM variables",
		Long:  `Plots heat maps of SOM variables`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			somFile := args[0]
			outFile := args[1]

			somYaml, err := os.ReadFile(somFile)
			if err != nil {
				return err
			}
			config, err := yml.ToSomConfig(somYaml)
			if err != nil {
				return err
			}

			s, err := som.New(config)
			if err != nil {
				return err
			}

			img, err := plot.Heatmap(&s, 0, 0, 250, 200)
			if err != nil {
				return err
			}

			err = os.MkdirAll(path.Dir(outFile), os.ModePerm)
			if err != nil {
				return err
			}
			file, err := os.Create(outFile)
			if err != nil {
				return err
			}

			return png.Encode(file, img)
		},
	}

	return command
}
