package main

import (
	"log"

	"github.com/mlange-42/som/cmd/som/cli"
)

func main() {
	if err := cli.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
