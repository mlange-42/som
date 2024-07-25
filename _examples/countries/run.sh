#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

som train untrained.yml data.csv -p out/progress.csv -P 10 > out/trained.yml
som export out/trained.yml > out/export.csv
som bmu out/trained.yml data.csv --ignore continent --preserve Country,code,continent > out/bmu.csv
som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --label Country --ignore continent
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --label Country --ignore continent
som plot density out/trained.yml out/density.png --data-file data.csv --label Country --ignore continent
som plot error out/trained.yml out/error.png --data-file data.csv --label Country --ignore continent
som plot xy out/trained.yml out/xy.png -x log_GNI -y Income_low_40 -c continent --data-file data.csv -C continent

som fill out/trained.yml data.csv --ignore continent --preserve Country,code,continent > out/filled.csv

som plot codes rose out/trained.yml out/codes-rose.png -n -s 800,600
