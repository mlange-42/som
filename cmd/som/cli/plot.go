package cli

import (
	"github.com/spf13/cobra"
)

func plotCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "plot [flags] <som-file>",
		Short: "Plots a SOM in various ways, see sub-commands",
		Long:  `Plots a SOM in various ways, see sub-commands`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	command.AddCommand(heatmapCommand())
	command.AddCommand(uMatrixCommand())
	command.AddCommand(xyCommand())
	command.AddCommand(densityCommand())

	return command
}
