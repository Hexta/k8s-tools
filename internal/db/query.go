package db

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/db/duckdb"
)

func Query(ctx context.Context, dataDir string, q string) (string, error) {
	return duckdb.Query(ctx, dataDir, q)
}
