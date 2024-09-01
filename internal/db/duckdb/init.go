package duckdb

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hexta/k8s-tools/internal/nodeutil"
	"github.com/marcboeker/go-duckdb"
	log "github.com/sirupsen/logrus"
)

const duckdbDir = "duckdb"

var (
	//go:embed sql/init.sql
	initSQL embed.FS
)

func InitDB(ctx context.Context, dataDir string, nodes []nodeutil.NodeInfo) error {
	dbDir := filepath.Join(dataDir, duckdbDir)

	if _, err := os.Stat(dbDir); !os.IsNotExist(err) {
		err := os.RemoveAll(dbDir)
		if err != nil {
			log.Errorf("failed to remove DB dir: %s", err)
		}
	}

	err := os.MkdirAll(dbDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create DB directory %s: %w", dbDir, err)
	}

	db, err := initConnection(dataDir)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Errorf("failed to close DB connection: %s", err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	err = runDDL(ctx, db)
	if err != nil {
		return err
	}

	return InsertNodes(ctx, db, nodes)
}

func initConnection(dataDir string) (db *sql.DB, err error) {
	dbDir := filepath.Join(dataDir, duckdbDir)
	dbFile := filepath.Join(dbDir, "duckdb.db")

	connector, err := duckdb.NewConnector(dbFile, nil)
	if err != nil {
		err = fmt.Errorf("failed to create DB connector: %w", err)
		return
	}

	db = sql.OpenDB(connector)

	return
}

func runDDL(ctx context.Context, db *sql.DB) error {
	query, err := initSQL.ReadFile("sql/init.sql")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(query))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}
