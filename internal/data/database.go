package data

import (
	"database/sql"
	"fmt"
	"go-manage/cmd/config"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (*sql.DB, error) {
	var conn *sql.DB
	var connErr error

	if !exists(config.DBPath) {
		fmt.Println("DATABASE DOESN'T EXIST. CREATING....")
		conn, connErr = sql.Open(config.DBDriver, config.DBPath)
		if connErr != nil {
			return nil, connErr
		}
		if createTableErr := createTable(conn, config.CreateTableQuery); createTableErr != nil {
			return nil, createTableErr
		}
	} else {
		fmt.Println("DATABASE FOUND. USING DATABASE....")
		conn, connErr = sql.Open(config.DBDriver, config.DBPath)
		if connErr != nil {
			return nil, connErr
		}
	}

	return conn, nil
}

func exists(dbPath string) bool {
	_, err := os.Stat(dbPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func createTable(db *sql.DB, tableQuery string) error {
	_, createTableErr := db.Exec(tableQuery)
	if createTableErr != nil {
		return fmt.Errorf("error creating table. Error: %w", createTableErr)
	}
	return nil
}
