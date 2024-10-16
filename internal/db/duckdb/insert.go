package duckdb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/k8s/container"
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
	ContainersTable       = "containers"
	DSTable               = "ds"
	DeploymentsTable      = "deployments"
	HPATable              = "hpa"
	NodesTable            = "nodes"
	PodsTable             = "pods"
	STSTable              = "sts"
	Schema                = "k8s"
	TaintsTable           = "taints"
	containerListCapacity = 65536
)

func InsertNodes(con driver.Conn, db *sql.DB, items k8snode.InfoList) error {
	return doInsert[k8snode.Info](con, db, Schema, NodesTable, items)
}

func InsertNodeTaints(con driver.Conn, db *sql.DB, items k8s.TaintList) error {
	return doInsert[k8s.Taint](con, db, Schema, TaintsTable, items)
}

func InsertPods(con driver.Conn, db *sql.DB, items k8spod.InfoList) error {
	return doInsert[k8spod.Info](con, db, Schema, PodsTable, items)
}

func InsertContainers(con driver.Conn, db *sql.DB, items k8spod.InfoList) error {
	containers := make(container.InfoList, 0, containerListCapacity)

	for _, pod := range items {
		containers = append(containers, pod.Containers...)
	}

	err := doInsert[container.Info](con, db, Schema, ContainersTable, containers)
	if err != nil {
		return err
	}

	return nil
}

func InsertDeployments(con driver.Conn, db *sql.DB, items deployment.InfoList) error {
	return doInsert[deployment.Info](con, db, Schema, DeploymentsTable, items)
}

func InsertHPAs(con driver.Conn, db *sql.DB, items hpa.InfoList) error {
	return doInsert[hpa.Info](con, db, Schema, HPATable, items)
}

func InsertSTS(con driver.Conn, db *sql.DB, items sts.InfoList) error {
	return doInsert[sts.Info](con, db, Schema, STSTable, items)
}

func InsertDS(con driver.Conn, db *sql.DB, items ds.InfoList) error {
	return doInsert[ds.Info](con, db, Schema, DSTable, items)
}

func doInsert[T any](con driver.Conn, db *sql.DB, schema string, table string, items []*T) error {
	columnNames, err := listTableColumnNames(db, schema, table)
	if err != nil {
		return err
	}

	columnIndexByName := getColumnIndexByName(columnNames)

	appender, err := duckdb.NewAppenderFromConn(con, schema, table)
	if err != nil {
		return err
	}

	for _, item := range items {
		rowValues, err := prepareRowValueSlice(*item, columnIndexByName)
		if err != nil {
			return err
		}

		err = appender.AppendRow(rowValues...)
		if err != nil {
			return err
		}
	}

	return appender.Flush()
}

func mapStringStringToDuckdbMap(m map[string]string) duckdb.Map {
	dm := make(duckdb.Map, len(m))

	for k, v := range m {
		dm[k] = v
	}

	return dm
}

func getColumnIndexByName(columns []string) map[string]int {
	columnIndexByName := make(map[string]int, len(columns))
	for i, name := range columns {
		columnIndexByName[name] = i
	}
	return columnIndexByName
}

func prepareRowValueSlice(item any, columnIndexByName map[string]int) ([]driver.Value, error) {
	values := make([]driver.Value, len(columnIndexByName))

	st := reflect.TypeOf(item)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tagValue := field.Tag.Get("db")
		if tagValue == "" {
			continue
		}
		fieldValue := reflect.ValueOf(item).Field(i).Interface()
		columnIndex, ok := columnIndexByName[tagValue]
		if !ok {
			return nil, fmt.Errorf("column %s not found", tagValue)
		}

		switch fieldValueTyped := fieldValue.(type) {
		case map[string]string:
			fieldValue = mapStringStringToDuckdbMap(fieldValueTyped)
		case *int32:
			fieldValue = *fieldValueTyped
		default:
		}

		values[columnIndex] = fieldValue
	}

	return values, nil
}
