package db

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/db/duckdb"
	"github.com/Hexta/k8s-tools/internal/nodeutil"
	"github.com/Hexta/k8s-tools/internal/podutil"
	log "github.com/sirupsen/logrus"
)

func InitDB(ctx context.Context, dataDir string, nodes []nodeutil.NodeInfo, pods []podutil.PodInfo) error {
	log.Debugf("Initializing database - start")
	defer log.Debugf("Initializing database - done")

	return duckdb.InitDB(ctx, dataDir, nodes, pods)
}
