package db

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/db/duckdb"
)

func RunTUI(ctx context.Context, dataDir string) {
	duckdb.RunTUI(ctx, dataDir)
}
