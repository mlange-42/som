package csv

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/layer"
)

func SomToCsv(som *som.Som, writer io.Writer, delim rune, noData string) error {
	layers := collectLayers(som)
	labelColumns, labels := collectLabels(som)

	err := writeHeadersSom(writer, labelColumns, layers, delim)
	if err != nil {
		return err
	}

	builder := strings.Builder{}

	nodes := som.Size().Nodes()
	for i := 0; i < nodes; i++ {
		x, y := som.Size().Coords(i)
		builder.WriteString(fmt.Sprintf("%d%c%d%c%d%c", i, delim, x, delim, y, delim))

		for j := range labels {
			builder.WriteString(labels[j][i])
			if i < len(labels)-1 || len(layers) > 0 {
				builder.WriteRune(delim)
			}
		}

		for j, layer := range layers {
			row := layer.GetNodeAt(i)
			for k, v := range row {
				if math.IsNaN(v) {
					builder.WriteString(noData)
				} else {
					builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
				}
				if k < len(row)-1 || j < len(layers)-1 {
					builder.WriteRune(delim)
				}
			}
		}

		if i < nodes-1 {
			builder.WriteRune('\n')
		}
		_, err := writer.Write([]byte(builder.String()))
		if err != nil {
			return err
		}
		builder.Reset()
	}
	return nil
}

func writeHeadersSom(writer io.Writer, labelColumns []string, layers []*layer.Layer, delim rune) error {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("node_id%cnode_x%cnode_y%c", delim, delim, delim))
	for i, col := range labelColumns {
		builder.WriteString(col)
		if i < len(labelColumns)-1 || len(layers) > 0 {
			builder.WriteRune(delim)
		}
	}

	for i, layer := range layers {
		cols := layer.ColumnNames()
		for j, col := range cols {
			builder.WriteString(col)
			if j < len(cols)-1 || i < len(layers)-1 {
				builder.WriteRune(delim)
			}
		}
	}

	builder.WriteRune('\n')
	_, err := writer.Write([]byte(builder.String()))
	return err
}

func collectLayers(som *som.Som) []*layer.Layer {
	layers := []*layer.Layer{}
	for _, lay := range som.Layers() {
		if lay.IsCategorical() {
			continue
		}

		lay, err := layer.NewWithData(
			lay.Name(), lay.ColumnNames(), lay.Normalizers(), *som.Size(),
			lay.Metric(), lay.Weight(), lay.IsCategorical(), append([]float64{}, lay.Weights()...))
		if err != nil {
			panic(err)
		}

		lay.DeNormalize()
		layers = append(layers, lay)
	}

	return layers
}

func collectLabels(som *som.Som) ([]string, [][]string) {
	labelColumns := []string{}
	labels := [][]string{}

	for _, layer := range som.Layers() {
		if !layer.IsCategorical() {
			continue
		}
		classes, indices := conv.LayerToClasses(layer)
		labs := make([]string, len(indices))
		for i := range indices {
			labs[i] = classes[indices[i]]
		}
		labels = append(labels, labs)
		labelColumns = append(labelColumns, layer.Name())
	}

	return labelColumns, labels
}
