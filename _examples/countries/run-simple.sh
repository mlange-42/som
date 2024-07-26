#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

echo Train SOM on countries dataset
som train untrained-simple.yml data.csv -p out/progress.csv -P 10 > out/trained-simple.yml
som quality out/trained-simple.yml data.csv

som plot heatmap out/trained-simple.yml out/heatmap-simple.png --data-file data.csv --label Country
som plot u-matrix out/trained-simple.yml out/u-matrix-simple.png --data-file data.csv --label Country
som plot xy out/trained-simple.yml out/xy-simple.png -x log_GNI -y Income_low_40 --data-file data.csv
som plot codes rose out/trained.yml out/codes-rose-simple.png -n -s 800,600

echo See sub-folder 'out/' for result tables and images.