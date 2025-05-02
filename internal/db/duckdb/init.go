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
	err := resetDBDirectory(dbDir)
	if err != nil {
		return err
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
			return InsertContainers(con, db, k8sInfo.Pods)
		},
		func() error {
			return InsertDeployments(con, db, k8sInfo.Deployments)
		},
		func() error {
			return InsertDSs(con, db, k8sInfo.DSs)
		},
		func() error {
			return InsertEndpoints(con, db, k8sInfo.EndpointSlices)
		},
		func() error {
			return InsertEndpointSlicePorts(con, db, k8sInfo.EndpointSlices)
		},
		func() error {
			return InsertEndpointSlices(con, db, k8sInfo.EndpointSlices)
		},
		func() error {
			return InsertHPAs(con, db, k8sInfo.HPAs)
		},
		func() error {
			return InsertImages(con, db, k8sInfo.Images)
		},
		func() error {
			return InsertInitContainers(con, db, k8sInfo.Pods)
		},
		func() error {
			return InsertNodes(con, db, k8sInfo.Nodes)
		},
		func() error {
			return InsertPods(con, db, k8sInfo.Pods)
		},
		func() error {
			return InsertPVs(con, db, k8sInfo.PVs)
		},
		func() error {
			return InsertServices(con, db, k8sInfo.Services)
		},
		func() error {
			return InsertSTS(con, db, k8sInfo.STSs)
		},
		func() error {
			return InsertTaints(con, db, k8sInfo.Taints)
		},
		func() error {
			return InsertTolerations(con, db, k8sInfo.Tolerations)
		},
	}

	for _, fn := range functionsToRun {
		if err = fn(); err != nil {
			return err
		}
	}

	return flush(ctx, db)
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

// resetDBDirectory ensures that the database directory is clean and ready.
func resetDBDirectory(dbDir string) error {
	// Clean up existing directory if it exists
	if _, err := os.Stat(dbDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(dbDir); err != nil {
			log.Errorf("Failed to remove DB dir: %s", err)
			return err
		}
	}

	// Create directory
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return fmt.Errorf("failed to create DB directory %s: %w", dbDir, err)
	}
	return nil
}

func flush(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "CHECKPOINT")
	return err
}
