package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB
var err error

func main() {
	db, err = getDbConnection("mysql", "abishek:mypassword^123@tcp(127.0.0.1:3306)/sample_db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	http.HandleFunc("/employees", EmployeeHandler)
	log.Println("Starting server at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
