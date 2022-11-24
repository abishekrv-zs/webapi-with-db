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

		var reqEmp []employee

		switch queryPram := r.URL.Query(); true {

		case queryPram.Has("id"):
			for _, emp := range allEmp {
				if emp.Id == queryPram.Get("id") {
					reqEmp = append(reqEmp, emp)
					break
				}
			}
			if reqEmp == nil {
				w.WriteHeader(http.StatusNotFound)
				if _, err := w.Write([]byte(`{"error":"emp id not found"}`)); err != nil {
					log.Println(err)
					return
				}
			}

		case queryPram.Has("deptid"):
			for _, emp := range allEmp {
				if emp.Dept.Id == queryPram.Get("deptid") {
					reqEmp = append(reqEmp, emp)
				}
			}
			if reqEmp == nil {
				w.WriteHeader(http.StatusNotFound)
				if _, err := w.Write([]byte(`{"error":"deptid not found"}`)); err != nil {
					log.Println(err)
				}
			}

		default:
			reqEmp = allEmp
		}

		// Marshal the required employees and send as response
		respBody, err := json.Marshal(reqEmp)
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
