package duckdb

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/Hexta/k8s-tools/internal/format"
	log "github.com/sirupsen/logrus"
)

func Query(ctx context.Context, dataDir string, q string) (*format.Data, error) {
	connector, err := initConnector(dataDir)

	if err != nil {
		return nil, err
	}

	db, _, err := initConnection(ctx, connector)

	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Errorf("failed to close rows: %v", err)
		}
	}(rows)

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	output := format.NewData()
	for _, name := range columnNames {
		output.AddColumn(name)
	}

	cr, err := newColumnsReader(rows)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		columns, err := cr.ReadColumns()
		if err != nil {
			return nil, err
		}
		output.AddRow(columns)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
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

	columnCount := len(columnNames)
	columns := make([]interface{}, columnCount)
	columnPtrs := make([]interface{}, columnCount)

	for i := range columns {
		columnPtrs[i] = &columns[i]
	}

	return &columnsReader{
		rows:       rows,
		columns:    columns,
		columnPtrs: columnPtrs,
	}, nil
}

func (cr *columnsReader) ReadColumns() ([]interface{}, error) {
	err := cr.rows.Scan(cr.columnPtrs...)
	if err != nil {
		return nil, err
	}

	return slices.Clone(cr.columns), nil
}
