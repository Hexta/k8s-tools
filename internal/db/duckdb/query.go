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
	connector, err := initConnector(dataDir)
	if err != nil {
		return "", err
	}

	con, _, err := initConnection(ctx, connector)

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

	var buf bytes.Buffer

	columnNamesInterface := make([]interface{}, len(columnNames))
	for i, name := range columnNames {
		columnNamesInterface[i] = name
	}

	tbl := table.New(columnNamesInterface...).WithWriter(&buf)
	cr, err := newColumnsReader(rows)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		columns, err := cr.ReadColumns()
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

type columnsReader struct {
	rows       *sql.Rows
	columns    []interface{}
	columnPtrs []interface{}
}

func newColumnsReader(rows *sql.Rows) (*columnsReader, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	columns := make([]interface{}, len(columnNames))
	columnPtrs := make([]interface{}, len(columnNames))

	for i := range columns {
		columnPtrs[i] = &columns[i]
	}

	return &columnsReader{
		rows:       rows,
		columns:    columns,
		columnPtrs: columnPtrs,
	}, nil
}

func (cr columnsReader) ReadColumns() ([]interface{}, error) {
	err := cr.rows.Scan(cr.columnPtrs...)
	if err != nil {
		return nil, err
	}

	return cr.columns, nil
}
