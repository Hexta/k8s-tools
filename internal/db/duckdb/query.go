package duckdb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/rodaine/table"
	log "github.com/sirupsen/logrus"
)

func Query(ctx context.Context, dataDir string, q string) (string, error) {
	con, err := initConnection(dataDir)

	if err != nil {
		return "", err
	}

	rows, err := con.QueryContext(ctx, q)
	if err != nil {
		return "", err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Errorf("failed to close rows: %v", err)
		}
	}(rows)

	columnNames, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("failed to get columns: %w", err)
	}

	columns := make([]interface{}, len(columnNames))
	columnPtrs := make([]interface{}, len(columnNames))

	for i := range columns {
		columnPtrs[i] = &columns[i]
	}

	var buf bytes.Buffer

	columnNamesInterface := make([]interface{}, len(columnNames))
	for i, name := range columnNames {
		columnNamesInterface[i] = name
	}

	tbl := table.New(columnNamesInterface...).WithWriter(&buf)

	for rows.Next() {
		err := rows.Scan(columnPtrs...)
		if err != nil {
			return "", err
		}
		tbl.AddRow(columns...)
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	tbl.Print()

	return buf.String(), nil
}
