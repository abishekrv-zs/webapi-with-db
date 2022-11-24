package main

import (
	"database/sql"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/employee", getEmployeeHandler).Methods("GET")
	router.HandleFunc("/employee", postEmployeeHandler).Methods("POST")
	router.HandleFunc("/department", getDepartmentHandler).Methods("GET")
	router.HandleFunc("/department", postDepartmentHandler).Methods("POST")

	//router.HandleFunc("/department{id}", getDepartmentByIdHandler).Methods("GET")
	//router.HandleFunc("/employee/{id}", getEmployeeByIdHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
