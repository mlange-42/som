package csv

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mlange-42/som"
	"github.com/mlange-42/som/conv"
	"github.com/mlange-42/som/layer"
)

func SomToCsv(som *som.Som, writer io.Writer, delim rune) error {
	labels := [][]string{}
	labelColumns := []string{}
	layers := []*layer.Layer{}

	for _, layer := range som.Layers() {
		if layer.IsCategorical() {
			classes, indices := conv.LayerToClasses(&layer)
			labs := make([]string, len(indices))
			for i := range indices {
				labs[i] = classes[indices[i]]
			}
			labels = append(labels, labs)
			labelColumns = append(labelColumns, layer.Name())
			continue
		}

		layers = append(layers, &layer)
	}

	del := string(delim)
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("id%snode_x%snode_y%s", del, del, del))

	for i, col := range labelColumns {
		builder.WriteString(col)
		if i < len(labelColumns)-1 || len(layers) > 0 {
			builder.WriteString(del)
		}
	}

	for i, layer := range layers {
		cols := layer.ColumnNames()
		for j, col := range cols {
			builder.WriteString(col)
			if j < len(cols)-1 || i < len(layers)-1 {
				builder.WriteString(del)
			}
		}
	}

	builder.WriteRune('\n')
	_, err := writer.Write([]byte(builder.String()))
	if err != nil {
		return err
	}
	builder.Reset()

	nodes := som.Size().Width * som.Size().Height
	for i := 0; i < nodes; i++ {
		x, y := som.Size().CoordsAt(i)
		builder.WriteString(fmt.Sprintf("%d%s%d%s%d%s", i, del, x, del, y, del))

		for j := range labels {
			builder.WriteString(labels[j][i])
			if i < len(labels)-1 || len(layers) > 0 {
				builder.WriteString(del)
			}
		}

		for j, layer := range layers {
			row := layer.GetNodeAt(i)
			for k, v := range row {
				builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
				if k < len(row)-1 || j < len(layers)-1 {
					builder.WriteString(del)
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
