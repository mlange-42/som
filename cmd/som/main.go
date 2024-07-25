package main

import (
	"fmt"
	"os"

	"github.com/mlange-42/som/cmd/som/cli"
)

func main() {
	command, err := cli.RootCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
