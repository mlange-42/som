# SOM -- Self-organizing maps in Go

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/som/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/som/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/som)](https://goreportcard.com/report/github.com/mlange-42/som)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/som)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/som)

[Go](https://go.dev) implementation of Self-Organizing Maps (SOM) alias Kohonen maps.
Provides a command line tool and a library for training and visualizing SOMs.

:warning: This is early work in progress!

## Features

* Multi-layered SOMs, alias XYF, alias super-SOMs.
* Supports continuous and discrete data.
* Training from CSV files without any manual preprocessing.
* Fully customizable training parameters.
* Visualization of SOMs through heatmaps with data point labels.
* Use as command line tool or as [Go](https://go.dev) library.

> Please note that the **built-in visualizations** are not intended for publication-quality output.
> Instead, they serve as quick tools for inspecting training and prediction results.
> For high-quality visualizations, we recommend exporting the SOM and other results to CSV files.
> You can then use dedicated visualization libraries in languages such as
> [Python](https://www.python.org/) or [R](https://www.r-project.org/) to create more refined and customized graphics.

## Installation

As long as there are no official releases, you can install the latest version from GitHub with Go:

```shell
go install github.com/mlange-42/som/cmd/som@latest
```

## Usage

Here is an example of how to use the command line tool, using the well-known World Countries dataset.

First, train a SOM with the Iris dataset:

```shell
som train _examples/countries/untrained.yml _examples/countries/data.csv > trained.yml
```

Visualize the trained SOM, showing labels of data points:

```shell
som plot heatmap trained.yml heatmap.png --data-file _examples/countries/data.csv --labels Country
```

You can also export the trained SOM to a CSV file:

```shell
som export trained.yml > nodes.csv
```

You can also determine the best-matching unit (BMU) for a each row in the dataset:

```shell
som bmu trained.yml _examples/countries/data.csv --preserve Country,code,continent > bmu.csv
```

## License

This project is distributed under the [MIT license](./LICENSE).
