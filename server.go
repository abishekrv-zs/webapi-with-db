package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func setDbConnection() (*sql.DB, error) {
	driver := "mysql"
	connectionString := "abishek:mypassword^123@tcp(127.0.0.1)/sample_db"

	db, err := sql.Open(driver, connectionString)

	return db, err
}
