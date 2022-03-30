package gormcallback

import (
	"database/sql"
	"fmt"

	"github.com/bndr/gotabulate"
	"gorm.io/gorm"
)

// RegisterExplainSQL register *gorm.DB instance with ExplainSQL callback,
// which makes the corresponding EXPLAIN of every SELECT, INSERT, UPDATE, DELETE statement printed to the stdout.
func RegisterExplainSQL(db *gorm.DB) error {
	return registerCallback(db, ExplainSQL)
}

// ExplainSQL makes the corresponding EXPLAIN of every SELECT, INSERT, UPDATE, DELETE statement printed to the stdout.
func ExplainSQL(db *gorm.DB) {
	stmt := db.Statement

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	explainSQL := fmt.Sprintf("EXPLAIN %s", stmt.SQL.String())
	rows, err := sqlDB.Query(explainSQL, stmt.Vars...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	results, err := tableDataFromSQLRows(rows)
	if err != nil {
		panic(err)
	}

	const alignRight = "right"
	const formatGrid = "grid"
	table := gotabulate.Create(results)
	table.SetHeaders(columns)
	table.SetAlign(alignRight)
	fmt.Println(stmt.SQL.String())
	fmt.Println(table.Render(formatGrid))
}

func tableDataFromSQLRows(rows *sql.Rows) ([][]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnCount := len(columns)
	var data [][]interface{}
	for rows.Next() {
		rowRawData := make([]sql.RawBytes, columnCount)
		scanDest := make([]interface{}, columnCount)
		for i := range scanDest {
			scanDest[i] = &rowRawData[i]
		}
		if serr := rows.Scan(scanDest...); serr != nil {
			return nil, serr
		}
		rowData := make([]interface{}, columnCount)
		for i := range rowRawData {
			rowData[i] = string(rowRawData[i])
		}
		data = append(data, rowData)
	}

	return data, nil
}
