# SOM -- Self-organizing maps in Go

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/som/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/som/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/som)](https://goreportcard.com/report/github.com/mlange-42/som)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/som)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/som)

Go implementation of Self-Organizing Maps (SOM) alias Kohonen maps.
Provides a command line tool and a library for training and visualizing SOMs.

:warning: This is early work in progress!

## Installation

As long as there are no official releases, you can install the latest version from GitHub with Go:

```shell
go install github.com/mlange-42/som/cmd/som@latest
```

## Usage

Here is an example of how to use the command line tool, using the well-known Iris dataset.

First, train a SOM with the Iris dataset:

```shell
som train _examples/iris/untrained.yml _examples/iris/data.csv > trained.yml
```

You can then export the trained SOM to a CSV file:

```shell
som export trained.yml > nodes.csv
```

You can also determine the best-matching unit (BMU) for a each row in the dataset:

```shell
som bmu trained.yml _examples/iris/data.csv --preserve species > bmu.csv
```
