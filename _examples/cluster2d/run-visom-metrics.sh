#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

echo Train and plot using ViSOM manhattan
som train untrained.yml data.csv -V manhattan > out/trained-vi-manhattan.yml
som plot xy out/trained-vi-manhattan.yml out/xy-vi-manhattan.png -x x -y y --data-file data.csv -C class

echo Train and plot using ViSOM euclidean
som train untrained.yml data.csv -V euclidean > out/trained-vi-euclidean.yml
som plot xy out/trained-vi-euclidean.yml out/xy-vi-euclidean.png -x x -y y --data-file data.csv -C class

echo Train and plot using ViSOM chebyshev
som train untrained.yml data.csv -V chebyshev > out/trained-vi-chebyshev.yml
som plot xy out/trained-vi-chebyshev.yml out/xy-vi-chebyshev.png -x x -y y --data-file data.csv -C class

echo See sub-folder 'out/' for result tables and images.
