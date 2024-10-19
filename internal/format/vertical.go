package format

import (
	"fmt"
	"strings"

	"github.com/rodaine/table"
)

func Vertical(input *Data) string {
	rows := input.GetRows()
	colNames := input.GetColumns()

	var sb strings.Builder

	for i, row := range rows {
		_, _ = fmt.Fprintf(&sb, "Row %d:\n%s\n", i, strings.Repeat("-", 7))
		tbl := table.New("", "").WithWriter(&sb).WithPrintHeaders(false)

		for j, col := range row {
			tbl.AddRow(colNames[j], col)
		}

		tbl.Print()
		sb.WriteString("\n")
	}

	return strings.TrimRight(sb.String(), "\n")
}
