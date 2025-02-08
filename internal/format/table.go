package format

import (
	"bytes"
	"strings"

	"github.com/rodaine/table"
)

func Table(input *Data, options Options) string {
	var buf bytes.Buffer
	var tbl table.Table

	if options.NoHeaders {
		// Create a table without headers.
		tbl = table.New(input.columns...).WithWriter(&buf).WithPrintHeaders(false)
	} else {
		// Create a table with headers and set a separator row.
		tbl = table.New(input.columns...).WithWriter(&buf).WithHeaderSeparatorRow('_')
	}

	for _, row := range input.rows {
		tbl.AddRow(row...)
	}

	tbl.Print()
	return strings.TrimRight(buf.String(), "\n")
}
