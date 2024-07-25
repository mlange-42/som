package cli

import (
	"fmt"

	"github.com/mlange-42/som/cmd/som/tree"
	"github.com/spf13/cobra"
)

// RootCommand sets up the CLI
func RootCommand() (*cobra.Command, error) {
	cobra.EnableCommandSorting = false

	root := &cobra.Command{
		Use:   "som",
		Short: "Self-organizing maps command line tool.",
		Long: `Self-organizing maps command line tool.

SOM is a tool for training, applying and analyzing self-organizing maps.

| Self-organizing maps are a type of artificial neural network that can be used
| for dimensionality reduction, classification, prediction and visualization of
| high-dimensional data.

Formats
=======

SOM definitions as well as trained SOMs are stored in YAML format (.yml).

CSV files are used as data input and output format (.csv).
Column delimiter and no-data strings can be configured using CLI flags.


For details, see the sub-commands.
Use 'som [command] --help' for more information about a command.`,
		SilenceUsage: true,
		//SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(trainCommand())
	root.AddCommand(qualityCommand())
	root.AddCommand(labelCommand())
	root.AddCommand(exportCommand())
	root.AddCommand(predictCommand())
	root.AddCommand(bmuCommand())
	root.AddCommand(fillCommand())
	root.AddCommand(plotCommand())

	addTreeToHelp(root, true)

	return root, nil
}

func addTreeToHelp(cmd *cobra.Command, setErrorPrefix bool) {
	cmdTree, err := tree.NewCmdTree(cmd)
	if err != nil {
		panic(err)
	}
	t, err := tree.FormatCmdTree(cmdTree, 2)
	if err != nil {
		panic(err)
	}
	cmd.Long += "\n\nCommand tree:\n\n" + t

	if setErrorPrefix {
		for _, name := range cmdTree.Nodes.Keys() {
			node, _ := cmdTree.Nodes.Get(name)
			node.Value.SetErrPrefix(fmt.Sprintf("Error in command '%s':", node.Value.CommandPath()))
		}
	}
}
