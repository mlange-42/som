mkdir out
som train untrained.yml data.csv > out/trained.yml
som export out/trained.yml > out/export.csv
som bmu out/trained.yml data.csv -p Country,code,continent > out/bmu.csv
som plot heatmap out/trained.yml out/heatmap.png --data-file data.csv --labels Country
som plot u-matrix out/trained.yml out/u-matrix.png --data-file data.csv --labels Country
