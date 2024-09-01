package duckdb

import (
	"context"
	"database/sql"

	"github.com/Hexta/k8s-tools/internal/nodeutil"
)

func InsertNodes(ctx context.Context, db *sql.DB, nodes []nodeutil.NodeInfo) error {
	for _, node := range nodes {
		_, err := db.ExecContext(
			ctx,
			"INSERT INTO k8s.nodes (name, creation_ts, instance_type, cpu_utilisation, memory_utilisation) VALUES (?, ?, ?, ?, ?)",
			node.Name,
			node.CreationTimestamp,
			node.InstanceType,
			node.Utilisation.CPU,
			node.Utilisation.Memory,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
