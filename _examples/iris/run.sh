#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

som train untrained.yml data.csv > out/trained.yml
som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --label species
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --label species
som plot xy out/trained.yml out/xy.png -x sepal_width -y petal_length -c species --data-file data.csv -C species

som plot codes line out/trained.yml out/codes-lines.png
som plot codes pie out/trained.yml out/codes-pie.png
som plot codes rose out/trained.yml out/codes-rose.png -z
