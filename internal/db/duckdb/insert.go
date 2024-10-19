package duckdb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"

	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/k8s/container"
	"github.com/Hexta/k8s-tools/internal/k8s/deployment"
	_ "github.com/Hexta/k8s-tools/internal/k8s/deployment"
	"github.com/Hexta/k8s-tools/internal/k8s/ds"
	"github.com/Hexta/k8s-tools/internal/k8s/hpa"
	k8snode "github.com/Hexta/k8s-tools/internal/k8s/node"
	k8spod "github.com/Hexta/k8s-tools/internal/k8s/pod"
	"github.com/Hexta/k8s-tools/internal/k8s/service"
	"github.com/Hexta/k8s-tools/internal/k8s/sts"
	"github.com/marcboeker/go-duckdb"
)

// Table names
const (
	ContainersTable     = "containers"
	DSTable             = "ds"
	DeploymentsTable    = "deployments"
	HPATable            = "hpa"
	InitContainersTable = "init_containers"
	NodesTable          = "nodes"
	PodsTable           = "pods"
	STSTable            = "sts"
	Schema              = "k8s"
	ServiceTable        = "services"
	TaintsTable         = "taints"
	TolerationsTable    = "tolerations"
)

const (
	containerListCapacity = 65536
)

func InsertNodes(con driver.Conn, db *sql.DB, items k8snode.InfoList) error {
	return doInsert[k8snode.Info](con, db, Schema, NodesTable, items)
}

func InsertTaints(con driver.Conn, db *sql.DB, items k8s.TaintList) error {
	return doInsert[k8s.Taint](con, db, Schema, TaintsTable, items)
}

func InsertTolerations(con driver.Conn, db *sql.DB, items k8s.TolerationList) error {
	return doInsert[k8s.Toleration](con, db, Schema, TolerationsTable, items)
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

func InsertInitContainers(con driver.Conn, db *sql.DB, items k8spod.InfoList) error {
	containers := make(container.InfoList, 0, containerListCapacity)

	for _, pod := range items {
		containers = append(containers, pod.InitContainers...)
	}

	err := doInsert[container.Info](con, db, Schema, InitContainersTable, containers)
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

func InsertDSs(con driver.Conn, db *sql.DB, items ds.InfoList) error {
	return doInsert[ds.Info](con, db, Schema, DSTable, items)
}

func InsertServices(con driver.Conn, db *sql.DB, items service.InfoList) error {
	return doInsert[service.Info](con, db, Schema, ServiceTable, items)
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

		fieldValue := reflect.ValueOf(item).Field(i)
		fieldValueInterface := fieldValue.Interface()

		columnIndex, ok := columnIndexByName[tagValue]
		if !ok {
			return nil, fmt.Errorf("column %s not found", tagValue)
		}

		switch fieldValueTyped := fieldValueInterface.(type) {
		case map[string]string:
			fieldValueInterface = mapStringStringToDuckdbMap(fieldValueTyped)
		case *bool:
			if fieldValueTyped == nil {
				fieldValueInterface = false
			} else {
				fieldValueInterface = *fieldValueTyped
			}
		case *int64, *int32:
			if fieldValue.IsNil() {
				fieldValueInterface = 0
			} else {
				elem := fieldValue.Elem()
				fieldValueInterface = elem.Interface()
			}
		case *string:
			if fieldValueTyped == nil {
				fieldValueInterface = ""
			} else {
				fieldValueInterface = *fieldValueTyped
			}
		default:
			if !slices.Contains([]string{"", "time"}, fieldValue.Type().PkgPath()) {
				if ok = fieldValue.CanConvert(reflect.TypeOf("")); ok {
					fieldValueInterface = fieldValue.Convert(reflect.TypeOf("")).Interface()
				}
			}
		}

		values[columnIndex] = fieldValueInterface
	}

	return values, nil
}
