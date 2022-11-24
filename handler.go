package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func getEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("select e.id, e.name, phone_number, department_id, d.name from employee e inner join department d on e.department_id = d.id")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
	}

	defer rows.Close()

	var allEmp []employee
	for rows.Next() {
		var emp employee
		if err := rows.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Dept.Id, &emp.Dept.Name); err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
			return
		}
		allEmp = append(allEmp, emp)
	}

	reqBody, err := json.Marshal(allEmp)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	w.Write(reqBody)
}

func getDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	rows, err := db.Query("select id, name from department")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	defer rows.Close()

	var allDep []department
	for rows.Next() {
		var dep department
		if err := rows.Scan(&dep.Id, &dep.Name); err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
			return
		}
		allDep = append(allDep, dep)
	}

	reqBody, err := json.Marshal(allDep)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	w.Write(reqBody)
}

func postEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}
	var emp employee
	if err := json.Unmarshal(reqBody, &emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	genUUID := strings.TrimSpace(string(uuid))

	_, err = db.Exec("insert into employee values(?,?,?,?)", genUUID, emp.Name, emp.Dept.Id, emp.PhoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	row := db.QueryRow("select e.id, e.name, phone_number, department_id, d.name from employee e inner join department d on e.department_id = d.id where e.id=?", genUUID)
	if err := row.Scan(&emp.Id, &emp.Name, &emp.PhoneNumber, &emp.Dept.Id, &emp.Dept.Name); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	respBody, err := json.Marshal(emp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

func postDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}
	var dep department
	if err := json.Unmarshal(reqBody, &dep); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	genUUID := strings.TrimSpace(string(uuid))

	_, err = db.Exec("insert into department values(?,?)", genUUID, dep.Name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	row := db.QueryRow("select id, name from department where id=?", genUUID)
	if err := row.Scan(&dep.Id, &dep.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}

	respBody, err := json.Marshal(dep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
