package format

import (
	"bytes"
	"strings"

	"github.com/rodaine/table"
)

func Table(input *Data) string {
	var buf bytes.Buffer
	tbl := table.New(input.columns...).WithWriter(&buf).WithHeaderSeparatorRow('_')

	for _, row := range input.rows {
		tbl.AddRow(row...)
	}

	tbl.Print()
	return strings.TrimRight(buf.String(), "\n")
}
