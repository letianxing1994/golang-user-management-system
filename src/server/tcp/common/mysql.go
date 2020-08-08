package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//build mysql connection
func NewMysqlConn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:811149@Tim@/user_db")
	if db == nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(10000)
	return
}

//get results, return one result
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}
