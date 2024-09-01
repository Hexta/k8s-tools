package db

import (
	"context"

	duckdb "github.com/Hexta/k8s-tools/internal/db/duckdb"
	"github.com/Hexta/k8s-tools/internal/nodeutil"
)

func InitDB(ctx context.Context, dataDir string, nodes []nodeutil.NodeInfo) error {
	return duckdb.InitDB(ctx, dataDir, nodes)
}
