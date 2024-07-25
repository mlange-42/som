#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

echo Train SOM on a dataset pf population pyramids of countries, and plot results
som train untrained.yml data.csv > out/trained.yml
som quality out/trained.yml data.csv

som plot heatmap out/trained.yml out/heatmap.png -f data.csv -l Country -s 600,360
som plot u-matrix out/trained.yml out/u-matrix.png -f data.csv -l Country -s 800,480
som plot density out/trained.yml out/density.png -f data.csv -l Country -s 800,480
som plot codes line out/trained.yml out/codes-line.png -s 800,480 --vertical --step pre --zero
som plot codes bar out/trained.yml out/codes-bar.png -s 800,480 --vertical --zero -C silver,white
som plot xy out/trained.yml out/xy-15-vs-60.png -x C15 -y C60 -f data.csv -s 800,600
