package db

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/db/duckdb"
	"github.com/Hexta/k8s-tools/internal/k8s"
	log "github.com/sirupsen/logrus"
)

func InitDB(ctx context.Context, dataDir string, k8sInfo *k8s.Info) error {
	log.Debugf("Initializing database - start")
	defer log.Debugf("Initializing database - done")

	return duckdb.InitDB(ctx, dataDir, k8sInfo)
}
