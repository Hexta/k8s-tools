package duckdb

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hexta/k8s-tools/internal/k8s"
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

	if err = runDDL(ctx, db); err != nil {
		return err
	}

	functionsToRun := []func() error{
		func() error {
			return InsertNodes(con, db, k8sInfo.Nodes)
		},
		func() error {
			return InsertNodeTaints(con, db, k8sInfo.NodeTaints)
		},
		func() error {
			return InsertPods(con, db, k8sInfo.Pods)
		},
		func() error {
			return InsertContainers(con, db, k8sInfo.Pods)
		},
		func() error {
			return InsertDeployments(con, db, k8sInfo.Deployments)
		},
		func() error {
			return InsertHPAs(con, db, k8sInfo.HPAs)
		},
		func() error {
			return InsertSTS(con, db, k8sInfo.STSs)
		},
		func() error {
			return InsertDSs(con, db, k8sInfo.DSs)
		},
		func() error {
			return InsertServices(con, db, k8sInfo.Services)
		},
	}

	for _, fn := range functionsToRun {
		if err = fn(); err != nil {
			return err
		}
	}

	return nil
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
