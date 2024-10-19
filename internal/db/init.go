package db

import (
	"context"
	"time"

	"github.com/Hexta/k8s-tools/internal/db/duckdb"
	"github.com/Hexta/k8s-tools/internal/k8s"
	log "github.com/sirupsen/logrus"
)

func InitDB(ctx context.Context, dataDir string, k8sInfo *k8s.Info) error {
	start := time.Now()
	log.Debugf("Initializing database")
	defer func() {
		log.Debugf("Initializing database: done, elapsed: %s", time.Since(start))
	}()

	return duckdb.InitDB(ctx, dataDir, k8sInfo)
}
