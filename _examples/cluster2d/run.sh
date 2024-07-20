mkdir out
som train untrained.yml data.csv -v 0.0 > out/trained.yml
som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --labels class
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --labels class
som plot xy out/trained.yml out/xy.png -x x -y y

som train untrained.yml data.csv > out/trained-vi.yml
som plot heatmap out/trained-vi.yml out/heatmap-vi.png --data-file data.csv --labels class
som plot u-matrix out/trained-vi.yml out/u-matrix-vi.png --data-file data.csv --labels class
som plot xy out/trained-vi.yml out/xy-vi.png -x x -y y