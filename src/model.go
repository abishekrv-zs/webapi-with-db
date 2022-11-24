package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func getDbConnection(driver string, connectionString string) (*sql.DB, error) {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

type department struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type employee struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	PhoneNumber string     `json:"phoneNumber"`
	Dept        department `json:"dept"`
}
