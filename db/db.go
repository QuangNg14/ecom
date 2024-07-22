package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN()) // Open a database specified by its database driver name and a driver-specific data source name
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
