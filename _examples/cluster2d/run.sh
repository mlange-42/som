#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

# Train and plot using simple SOM
som train untrained.yml data.csv -v 0.0 > out/trained.yml
som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --label class
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --label class
som plot xy out/trained.yml out/xy.png -x x -y y --data-file data.csv -C class

# Train and plot using ViSOM
som train untrained.yml data.csv > out/trained-vi.yml
som plot heatmap out/trained-vi.yml out/heatmap-vi.png --data-file data.csv --label class
som plot u-matrix out/trained-vi.yml out/u-matrix-vi.png --data-file data.csv --label class
som plot xy out/trained-vi.yml out/xy-vi.png -x x -y y --data-file data.csv -C class

# Perform label propagation on ViSOM, with very few labelled samples
som label out/trained-vi.yml data-labels.csv -c class > out/labelled.yml
som plot heatmap out/labelled.yml out/heatmap-labelled.png --data-file data-labels.csv --label class --ignore class -b class
som plot heatmap out/labelled.yml out/heatmap-labelled-all.png --data-file data.csv --label class --ignore class -b class
som plot xy out/labelled.yml out/xy-labelled.png -x x -y y -c class --data-file data.csv -C class --ignore class

# Predict labels on trained SOM
som predict out/labelled.yml data.csv --layers class --all > out/prediced.csv

som quality out/labelled.yml data.csv