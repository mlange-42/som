#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

echo Train SOM on digits dataset for character recognition
som train untrained.yml data.csv > out/trained.yml
som quality out/trained.yml data.csv

som plot heatmap out/trained.yml out/classes.png -c class -b class
som plot codes image out/trained.yml out/image.png -r 8 -s 900,600

echo See sub-folder 'out/' for result tables and images.