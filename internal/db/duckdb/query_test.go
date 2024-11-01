package duckdb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_columnsReader_ReadColumns_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Mock rows with data
	rows := sqlmock.NewRows([]string{"column1", "column2"}).
		AddRow(1, "test").
		AddRow(2, "example")

	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT * FROM test")
	require.NoError(t, err)

	cr, err := newColumnsReader(queryRows)
	require.NoError(t, err)

	expectedRows := [][]interface{}{{int64(1), "test"}, {int64(2), "example"}}
	actualRows := make([][]interface{}, 0)

	for queryRows.Next() {
		columns, err := cr.ReadColumns()
		require.NoError(t, err)
		actualRows = append(actualRows, columns)
	}

	assert.ElementsMatch(t, expectedRows, actualRows)
}
