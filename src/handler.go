package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func EmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":

		// Query the DB to fetch rows
		rows, err := db.Query("select e.id, e.name, phone_number, department_id,d.name from employee e inner join department d on e.department_id = d.id")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		allEmp := make([]employee, 0)
		for rows.Next() {
			var emp employee
			if err := rows.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Dept.Id, &emp.Dept.Name); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			allEmp = append(allEmp, emp)
		}

		// Marshal the required employees and send as response
		respBody, err := json.Marshal(allEmp)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(respBody); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}
