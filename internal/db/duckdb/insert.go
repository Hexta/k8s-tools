package duckdb

import (
	"context"
	"database/sql/driver"

	"github.com/Hexta/k8s-tools/internal/nodeutil"
	"github.com/marcboeker/go-duckdb"
)

func InsertNodes(_ context.Context, con driver.Conn, nodes []nodeutil.NodeInfo) error {
	appender, err := duckdb.NewAppenderFromConn(con, "k8s", "nodes")
	if err != nil {
		return err
	}

	for _, node := range nodes {
		err := appender.AppendRow(
			node.Name,
			node.CreationTimestamp,
			node.InstanceType,
			node.Utilisation.CPU,
			node.Utilisation.Memory,
			mapStringStringToDuckdbMap(node.Labels),
		)
		if err != nil {
			return err
		}
	}
	err = appender.Flush()
	if err != nil {
		return err
	}

	return nil
}

func mapStringStringToDuckdbMap(m map[string]string) duckdb.Map {
	dm := make(duckdb.Map, len(m))

	for k, v := range m {
		dm[k] = v
	}

	return dm
}
