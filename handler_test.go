package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllEmployees(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	h := myHandler{db: db}

	t.Run("Get all Employee", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/employees", nil)
		w := httptest.NewRecorder()

		mock.ExpectQuery(getAllEmployeeQuery).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phoneNumber", "department_id", "name"}).
				AddRow("uuid-for-emp", "testEmp", "9876543210", "uuid-for-dept", "testDept"))

		h.getAllEmployees(w, req)
		resp := w.Result()
		respBody, _ := io.ReadAll(resp.Body)

		expBody := `[{"id":"uuid-for-emp","name":"testEmp","phoneNumber":"9876543210","department":{"id":"uuid-for-dept","name":"testDept"}}]`
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expBody, string(respBody))
	})

	t.Run("Get all employee by a valid dept_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/employees?dept_id=valid-dept-id", nil)
		w := httptest.NewRecorder()

		mock.ExpectQuery(getAllEmployeeByDeptIdQuery).WithArgs("valid-dept-id").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phoneNumber", "department_id", "name"}).
				AddRow("uuid-for-emp", "testEmp", "9876543210", "valid-dept-id", "testDept"))

		h.getAllEmployees(w, req)

		resp := w.Result()
		respBody, _ := io.ReadAll(resp.Body)

		expBody := `[{"id":"uuid-for-emp","name":"testEmp","phoneNumber":"9876543210","department":{"id":"valid-dept-id","name":"testDept"}}]`
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expBody, string(respBody))
	})

	t.Run("Get all employee by invalid dept_id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/employees?dept_id=invalid-dept-id", nil)
		w := httptest.NewRecorder()

		mock.ExpectQuery(getAllEmployeeByDeptIdQuery).WithArgs("invalid-dept-id").
			WillReturnRows(sqlmock.NewRows([]string{}))

		h.getAllEmployees(w, req)
		resp := w.Result()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestGetEmployeeById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}

	h := myHandler{db: db}

	t.Run("Get employee by valid id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/employees/{id}", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "valid-emp-id"})
		w := httptest.NewRecorder()

		mock.ExpectQuery(getEmployeeByIdQuery).WithArgs("valid-emp-id").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phoneNumber", "department_id", "name"}).
				AddRow("valid-emp-id", "testEmp", "9876543210", "uuid-for-dept", "testDept"))

		h.getEmployeeById(w, req)

		resp := w.Result()
		respBody, _ := io.ReadAll(resp.Body)

		expResp := `{"id":"valid-emp-id","name":"testEmp","phoneNumber":"9876543210","department":{"id":"uuid-for-dept","name":"testDept"}}`

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expResp, string(respBody))
	})

	t.Run("Get employee by invalid id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/employees/{id}", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-emp-id"})
		w := httptest.NewRecorder()

		mock.ExpectQuery(getEmployeeByIdQuery).WithArgs("invalid-emp-id").
			WillReturnRows(sqlmock.NewRows([]string{}))

		h.getEmployeeById(w, req)

		resp := w.Result()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
