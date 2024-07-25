package main

import (
	"log"

	"github.com/mlange-42/som/cmd/som/cli"
)

func main() {
	command, err := cli.RootCommand()
	if err != nil {
		log.Fatal(err)
	}
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
