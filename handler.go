package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type myHandler struct {
	db *sql.DB
}

func (h myHandler) getAllEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if r.URL.Query().Has("dept_id") {
		deptId := r.URL.Query().Get("dept_id")

		rows, err := h.db.Query(getAllEmployeeByDeptIdQuery, deptId)
		if err != nil {
			log.Println(err)
			return
		}
		defer rows.Close()

		var allEmp []employee
		for rows.Next() {
			var emp employee
			if err := rows.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Department.Id, &emp.Department.Name); err != nil {
				log.Println(err)
				return
			}
			allEmp = append(allEmp, emp)
		}
		if len(allEmp) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		respBody, err := json.Marshal(allEmp)
		if _, err := w.Write(respBody); err != nil {
			log.Println(err)
			return
		}

	} else {
		rows, err := h.db.Query(getAllEmployeeQuery)
		if err != nil {
			log.Println(err)
			return
		}
		defer rows.Close()

		var allEmp []employee
		for rows.Next() {
			var emp employee
			if err := rows.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Department.Id, &emp.Department.Name); err != nil {
				log.Println(err)
				return
			}
			allEmp = append(allEmp, emp)
		}

		respBody, err := json.Marshal(allEmp)
		if _, err := w.Write(respBody); err != nil {
			log.Println(err)
			return
		}

	}
}

func (h myHandler) getEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	empId := mux.Vars(r)["id"]

	var emp employee
	row := h.db.QueryRow(getEmployeeByIdQuery, empId)
	if err := row.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Department.Id, &emp.Department.Name); err != nil {
		if err.Error() == "sql: no rows in result set" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.Println(err)
		return
	}

	respBody, err := json.Marshal(emp)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := w.Write(respBody); err != nil {
		log.Println(err)
		return
	}
}

func (h myHandler) postEmployee(w http.ResponseWriter, r *http.Request) {

}

func (h myHandler) getAllDepartments(w http.ResponseWriter, r *http.Request) {

}

func (h myHandler) getDepartmentById(w http.ResponseWriter, r *http.Request) {

}

func (h myHandler) postDepartment(w http.ResponseWriter, r *http.Request) {

}
