package duckdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/marcboeker/go-duckdb"
	log "github.com/sirupsen/logrus"
)

const duckdbDir = "duckdb"

var (
	//go:embed sql/init.sql
	initSQL embed.FS
)

func InitDB(ctx context.Context, dataDir string, k8sInfo *k8s.Info) error {
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

	connector, err := initConnector(dataDir)
	if err != nil {
		return fmt.Errorf("failed to initialize DB connector: %w", err)
	}

	db, con, err := initConnection(ctx, connector)
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

	err = InsertNodes(con, k8sInfo.Nodes)
	if err != nil {
		return err
	}

	err = InsertPods(con, k8sInfo.Pods)
	if err != nil {
		return err
	}

	err = InsertContainers(con, k8sInfo.Pods)
	if err != nil {
		return err
	}

	err = InsertDeployments(con, k8sInfo.Deployments)
	if err != nil {
		return err
	}

	err = InsertHPAs(con, k8sInfo.HPAs)
	if err != nil {
		return err
	}

	err = InsertSTS(con, k8sInfo.STSs)
	if err != nil {
		return err
	}

	return nil
}

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
