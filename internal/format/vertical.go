package format

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

func Vertical(input *Data) string {
	rows := input.GetRows()
	colNames := input.GetColumns()

	var sb strings.Builder

	for i, row := range rows {
		_, _ = fmt.Fprintf(&sb, "Row %d:\n%s\n", i+1, strings.Repeat("-", 7))
		tw := tabwriter.NewWriter(&sb, 0, 0, 1, ' ', 0)

		for j, col := range row {
			_, _ = fmt.Fprintf(tw, "%s\t%v\n", colNames[j], col)
		}

		_ = tw.Flush()
		sb.WriteString("\n")
	}

	return strings.TrimRight(sb.String(), "\n")
}
