#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

echo Train simple SOM on the Olive oils dataset and plot results
som train untrained.yml data.csv -v 0.0 > out/trained.yml
som quality out/trained.yml data.csv

som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --label id -b area
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --label id -b area
som plot xy out/trained.yml out/xy.png -x palmitic -y linoleic -c area --data-file data.csv -C region
som plot density out/trained.yml out/density.png --data-file data.csv --label id -b area
som plot error out/trained.yml out/error.png --data-file data.csv --label id -b area

echo Train ViSOM on the Olive oils dataset and plot results
som train untrained.yml data.csv > out/trained-vi.yml
som quality out/trained.yml data.csv

som plot heatmap out/trained-vi.yml out/heatmap-vi.png --data-file data.csv --label id -b area
som plot u-matrix out/trained-vi.yml out/u-matrix-vi.png --data-file data.csv --label id -b area
som plot xy out/trained-vi.yml out/xy-vi.png -x palmitic -y linoleic -c area --data-file data.csv -C region
som plot density out/trained-vi.yml out/density-vi.png --data-file data.csv --label id -b area
som plot error out/trained-vi.yml out/error-vi.png --data-file data.csv --label id -b area

echo See sub-folder 'out/' for result tables and images.
