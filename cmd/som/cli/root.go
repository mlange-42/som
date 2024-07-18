package cli

import (
	"github.com/spf13/cobra"
)

// RootCommand sets up the CLI
func RootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "som",
		Short:         "Self-organizing maps in Go",
		Long:          `Self-organizing maps in Go`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(trainCommand())

	return root
}
