# SOM -- Self-organizing maps in Go

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/som/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/som/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/som)](https://goreportcard.com/report/github.com/mlange-42/som)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/som)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/som)

[Go](https://go.dev) implementation of Self-Organizing Maps (SOM) alias Kohonen maps.
Provides a command line tool and a library for training and visualizing SOMs.

:warning: This is early work in progress!

![Heatmaps world countries dataset](https://github.com/user-attachments/assets/e01d4947-183c-4441-8a17-15f09d9f9e7e)  
*SOM heatmap visualization of the World Countries dataset.*

## Features

* Multi-layered SOMs, alias XYF, alias super-SOMs.
* Visualization Induced SOMs, alias ViSOMs.
* Training from CSV files without any manual preprocessing.
* Supports continuous and discrete data.
* Fully customizable training and SOM parameters.
* Visualization of SOMs by a wide range of flexible plots.
* Use as command line tool or as [Go](https://go.dev) library.

> Please note that the **built-in visualizations** are not intended for publication-quality output.
> Instead, they serve as quick tools for inspecting training and prediction results.
> For high-quality visualizations, we recommend exporting the SOM and other results to CSV files.
> You can then use dedicated visualization libraries in languages such as
> [Python](https://www.python.org/) or [R](https://www.r-project.org/) to create more refined and customized graphics.

## Installation

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/som/releases).

> Alternatively, install the latest version using [Go](https://go.dev):
> ```shell
> go install github.com/mlange-42/som/cmd/som@latest
> ```

## Usage

Get **help** for the command line tool:

```shell
som --help
```

Here are some examples how to use the command line tool, using the World Countries dataset.

**Train** an SOM with the dataset:

```shell
som train _examples/countries/untrained.yml _examples/countries/data.csv > trained.yml
```

**Visualize** the trained SOM as heatmaps of components, showing labels of data points (i.e. countries):

```shell
som plot heatmap trained.yml heatmap.png --data-file _examples/countries/data.csv --label Country
```

**Export** the trained SOM to a CSV file:

```shell
som export trained.yml > nodes.csv
```

Determine the **best-matching unit** (BMU) for a each row in the dataset:

```shell
som bmu trained.yml _examples/countries/data.csv --preserve Country,code,continent > bmu.csv
```

### Available commands

Taken from the CLI help, here is a tree representation of all currently available (sub)-commands:

```
som          Self-organizing maps command line tool.
├─train      Trains an SOM on the given dataset.
├─label      Classifies SOM nodes using label propagation.
├─export     Exports an SOM to a CSV table of node vectors.
├─predict    Predict entire layers or table columns using a trained SOM.
├─bmu        Finds the best-matching unit (BMU) for each table row in a dataset.
├─fill       Fills missing data in the data file based on a trained SOM.
└─plot       Plots visualizations for an SOM in various ways. See sub-commands.
  ├─heatmap  Plots heat maps of multiple SOM variables, a.k.a. components plot.
  ├─codes    Plots SOM node codes in different ways. See sub-commands.
  │ ├─line   Plots SOM node codes as line charts.
  │ ├─pie    Plots SOM node codes as pie charts.
  │ ├─rose   Plots SOM node codes as rose alias Nightingale charts.
  │ └─image  Plots SOM node codes as images.
  ├─u-matrix Plots the u-matrix of an SOM, showing inter-node distances.
  ├─xy       Plots for pairs of SOM variables as scatter plots.
  ├─density  Plots the data density of an SOM as a heatmap.
  └─error    Plots (root) mean-squared node error as a heatmap.
```

### YAML configuration

The command line tool uses a YAML configuration file to specify the SOM parameters.

Here is an example of a configuration file for the Iris dataset.
The dataset has these columns: `species`, `sepal_length`, `sepal_width`, `petal_length`, and `petal_width`.

```yaml
som:                     # SOM definitions
  size: [8, 6]           # Size of the SOM
  neighborhood: gaussian # Neighborhood function
  metric: manhattan      # Distance metric in map space

  layers:                # Layers of the SOM
    - name: Scalars      # Name of the layer. Has no meaning for continuous layers
      columns:           # Columns of the layer
        - sepal_length   # Column names as in the dataset
        - sepal_width
        - petal_length
        - petal_width
      norm: [gaussian]   # Normalization function(s) for columns
      metric: euclidean  # Distance metric
      weight: 1          # Weight of the layer

    - name: species      # Name of the layer. Use column name for categorical layers
      metric: hamming    # Distance metric
      categorical: true  # Layer is categorical. Omit columns
      weight: 0.5        # Weight of the layer

training:                # Training parameters. Optional. Can be overwritten by CLI arguments
  epochs: 2500                        # Number of training epochs
  alpha: polynomial 0.25 0.01 2       # Learning rate decay function
  radius: polynomial 6 1 2            # Neighborhood radius decay function
  weight-decay: polynomial 0.5 0.0 3  # Weight decay coefficient function
  lambda: 0.33                        # ViSOM resolution parameter
```

See the [examples](./_examples) folder for more examples.

## License

This project is distributed under the [MIT license](./LICENSE).
