#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

som train untrained.yml data.csv > out/trained.yml
som export out/trained.yml > out/export.csv
som bmu out/trained.yml data.csv -p Country,code,continent > out/bmu.csv
som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --labels Country
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --labels Country
som plot density out/trained.yml out/density.png --data-file data.csv --labels Country
som plot xy out/trained.yml out/xy.png -x log_GNI -y Income_low_40 -c continent --data-file data.csv -C continent
