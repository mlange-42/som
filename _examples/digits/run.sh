#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

som train untrained.yml data.csv > out/trained.yml
som plot heatmap out/trained.yml out/classes.png -c class
som plot codes image out/trained.yml out/image.png -r 8 -s 1200,800
