package cli

import (
	"github.com/mlange-42/som/cmd/som/tree"
	"github.com/spf13/cobra"
)

// RootCommand sets up the CLI
func RootCommand() (*cobra.Command, error) {
	cobra.EnableCommandSorting = false

	root := &cobra.Command{
		Use:           "som",
		Short:         "Self-organizing maps in Go.",
		Long:          `Self-organizing maps in Go.`,
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
	root.AddCommand(labelCommand())
	root.AddCommand(exportCommand())
	root.AddCommand(predictCommand())
	root.AddCommand(bmuCommand())
	root.AddCommand(fillCommand())
	root.AddCommand(plotCommand())

	t, err := tree.FormatCmdTree(root, 2)
	if err != nil {
		return nil, err
	}
	root.Long += "\n\nCommand tree:\n\n" + t

	return root, nil
}
