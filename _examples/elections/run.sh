#!/bin/bash
cd "$(dirname "$0")"

mkdir out
set -e

# Train and plot using simple SOM
som train untrained.yml data.csv -v 0.0 > out/trained.yml
som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --labels land
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --labels land
som plot xy out/trained.yml out/xy-cdu-spd.png -x CDU -y SPD --data-file data.csv -C land
som plot xy out/trained.yml out/xy-afd-gruene.png -x AfD -y GRUENE --data-file data.csv -C land

# Train and plot using ViSOM
som train untrained.yml data.csv > out/trained-vi.yml
som plot heatmap out/trained-vi.yml out/heatmap-vi.png --data-file data.csv --labels land
som plot u-matrix out/trained-vi.yml out/u-matrix-vi.png --data-file data.csv --labels land
som plot xy out/trained-vi.yml out/xy-cdu-spd-vi.png -x CDU -y SPD --data-file data.csv -C land
som plot xy out/trained-vi.yml out/xy-afd-gruene-vi.png -x AfD -y GRUENE --data-file data.csv -C land

# Plot codes as pie charts
som plot codes pie out/trained.yml out/codes.png -s 1200,800 \
        -c CDU,SPD,GRUENE,FDP,AfD,LINKE,PIRATEN,Other \
        -C black,red,green,yellow,blue,purple,orange,silver
som plot codes pie out/trained-vi.yml out/codes-vi.png -s 1200,800 \
        -c CDU,SPD,GRUENE,FDP,AfD,LINKE,PIRATEN,Other \
        -C black,red,green,yellow,blue,purple,orange,silver
