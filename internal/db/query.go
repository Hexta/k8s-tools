package db

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/db/duckdb"
	"github.com/Hexta/k8s-tools/internal/format"
)

func Query(ctx context.Context, dataDir string, q string) (*format.Data, error) {
	return duckdb.Query(ctx, dataDir, q)
}
