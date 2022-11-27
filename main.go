package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	db, err := setDbConnection()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	h := myHandler{db: db}

	router := mux.NewRouter()

	router.Handle("/", http.NotFoundHandler())

	router.HandleFunc("/employees", h.getAllEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", h.getEmployeeById).Methods("GET")
	router.HandleFunc("/employees", h.postEmployee).Methods("POST")

	router.HandleFunc("/departments", h.getAllDepartments).Methods("GET")
	router.HandleFunc("/departments/{id}", h.getDepartmentById).Methods("GET")
	router.HandleFunc("/departments", h.postDepartment).Methods("POST")

	log.Println("Starting server at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", router))
}
