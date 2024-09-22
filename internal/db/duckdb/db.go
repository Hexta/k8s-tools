package duckdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"path/filepath"

	"github.com/marcboeker/go-duckdb"
	log "github.com/sirupsen/logrus"
)

func initConnector(dataDir string) (*duckdb.Connector, error) {
	dbDir := filepath.Join(dataDir, duckdbDir)
	dbFile := filepath.Join(dbDir, "duckdb.db")

	return duckdb.NewConnector(dbFile, nil)
}

func initConnection(ctx context.Context, connector *duckdb.Connector) (db *sql.DB, connect driver.Conn, err error) {
	connect, err = connector.Connect(ctx)
	if err != nil {
		err = fmt.Errorf("failed to connect to DB: %w", err)
		return
	}

	db = sql.OpenDB(connector)

	return
}

func listTableColumnNames(db *sql.DB, schema string, table string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT column_name from (DESCRIBE TABLE %s.%s)", schema, table))
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("Failed to close statement: %v", err)
		}
	}()

	cols := make([]string, 0, 16)
	for rows.Next() {
		var col string
		err := rows.Scan(&col)
		if err != nil {
			return nil, err
		}

		cols = append(cols, col)
	}

	return cols, nil
}
