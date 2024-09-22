package duckdb

import (
	"database/sql/driver"

	"github.com/Hexta/k8s-tools/internal/k8s/deployment"
	_ "github.com/Hexta/k8s-tools/internal/k8s/deployment"
	"github.com/Hexta/k8s-tools/internal/k8s/ds"
	"github.com/Hexta/k8s-tools/internal/k8s/hpa"
	k8snode "github.com/Hexta/k8s-tools/internal/k8s/node"
	k8spod "github.com/Hexta/k8s-tools/internal/k8s/pod"
	"github.com/Hexta/k8s-tools/internal/k8s/sts"
	"github.com/marcboeker/go-duckdb"
)

const (
	Schema           = "k8s"
	NodesTable       = "nodes"
	PodsTable        = "pods"
	ContainersTable  = "containers"
	DeploymentsTable = "deployments"
	HPATable         = "hpa"
	STSTable         = "sts"
	DSTable          = "ds"
)

func InsertNodes(con driver.Conn, nodes k8snode.InfoList) error {
	appender, err := duckdb.NewAppenderFromConn(con, Schema, NodesTable)
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
			mapStringStringToDuckdbMap(node.Address),
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

func InsertPods(con driver.Conn, pods k8spod.InfoList) error {
	appender, err := duckdb.NewAppenderFromConn(con, Schema, PodsTable)
	if err != nil {
		return err
	}

	for _, pod := range pods {
		err := appender.AppendRow(
			pod.Name,
			pod.Namespace,
			pod.NodeName,
			pod.CreationTimestamp,
			mapStringStringToDuckdbMap(pod.Labels),
			pod.CPURequests,
			pod.CPULimits,
			pod.MemoryRequests,
			pod.MemoryLimits,
			pod.IP,
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

func InsertContainers(con driver.Conn, pods k8spod.InfoList) error {
	appender, err := duckdb.NewAppenderFromConn(con, Schema, ContainersTable)
	if err != nil {
		return err
	}

	for _, pod := range pods {
		for _, container := range pod.Containers {
			err := appender.AppendRow(
				container.Name,
				container.Namespace,
				container.PodName,
				container.CPURequests,
				container.CPULimits,
				container.MemoryRequests,
				container.MemoryLimits,
			)
			if err != nil {
				return err
			}
		}
	}

	err = appender.Flush()
	if err != nil {
		return err
	}

	return nil
}

func InsertDeployments(con driver.Conn, deployments deployment.InfoList) error {
	appender, err := duckdb.NewAppenderFromConn(con, Schema, DeploymentsTable)
	if err != nil {
		return err
	}

	for _, item := range deployments {
		replicas := int32(0)
		if item.Replicas != nil {
			replicas = *item.Replicas
		}

		err := appender.AppendRow(
			item.Name,
			item.Namespace,
			item.CreationTimestamp,
			mapStringStringToDuckdbMap(item.Labels),
			replicas,
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

func InsertHPAs(con driver.Conn, items hpa.InfoList) error {
	appender, err := duckdb.NewAppenderFromConn(con, Schema, HPATable)
	if err != nil {
		return err
	}

	for _, item := range items {
		err := appender.AppendRow(
			item.Name,
			item.Namespace,
			item.CreationTimestamp,
			mapStringStringToDuckdbMap(item.Labels),
			item.CurrentReplicas,
			item.DesiredReplicas,
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

func InsertSTS(con driver.Conn, items sts.InfoList) error {
	appender, err := duckdb.NewAppenderFromConn(con, Schema, STSTable)
	if err != nil {
		return err
	}

	for _, item := range items {
		err := appender.AppendRow(
			item.Name,
			item.Namespace,
			item.CreationTimestamp,
			mapStringStringToDuckdbMap(item.Labels),
			item.Replicas,
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

func InsertDS(con driver.Conn, items ds.InfoList) error {
	appender, err := duckdb.NewAppenderFromConn(con, Schema, DSTable)
	if err != nil {
		return err
	}

	for _, item := range items {
		err := appender.AppendRow(
			item.Name,
			item.Namespace,
			item.CreationTimestamp,
			mapStringStringToDuckdbMap(item.Labels),
			item.CurrentNumberScheduled,
			item.DesiredNumberScheduled,
			item.NumberAvailable,
			item.NumberMisscheduled,
			item.NumberReady,
			item.NumberUnavailable,
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
