package tree

import (
	"fmt"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

// FormatCmdTree creates a tree-like representation of a command and its sub-commands
func FormatCmdTree(cmdTree *CmdTree, indent int) (string, error) {
	wMax := 0
	for _, n := range cmdTree.Nodes.Keys() {
		node, _ := cmdTree.Nodes.Get(n)
		ln := utf8.RuneCountInString(node.Value.Command.Name()) + indent*node.Depth
		if ln > wMax {
			wMax = ln
		}
	}

	formatter := NewTreeFormatter(
		func(t *CmdNode, indent int) string {
			return fmt.Sprintf("%-*s %s", wMax-indent, t.Value.Command.Name(), t.Value.Short)
		},
		indent,
	)
	return formatter.FormatTree(cmdTree), nil
}

// CmdTree is a tree of cobra commands
type CmdTree = MapTree[CmdWrapper]

// CmdNode is a tree of cobra commands
type CmdNode = MapNode[CmdWrapper]

// CmdWrapper wraps *cobra.Command to implement the Named interface
type CmdWrapper struct {
	*cobra.Command
}

// GetName implements the Named interface required for the MapTree
func (cmd CmdWrapper) GetName() string {
	return nodePath(cmd.Command)
}

// NewCmdTree creates a new project tree
func NewCmdTree(command *cobra.Command) (*CmdTree, error) {

	t := NewTree(
		CmdWrapper{command},
	)

	err := buildTree(t, t.Root)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func nodePath(command *cobra.Command) string {
	if command.HasParent() {
		return fmt.Sprintf("%s/%s", nodePath(command.Parent()), command.Name())
	}
	return command.Name()
}

func buildTree(t *CmdTree, node *MapNode[CmdWrapper]) error {
	for _, cmd := range node.Value.Commands() {
		child, err := t.Add(node, CmdWrapper{cmd})
		if err != nil {
			return err
		}
		err = buildTree(t, child)
		if err != nil {
			return err
		}
	}
	return nil
}