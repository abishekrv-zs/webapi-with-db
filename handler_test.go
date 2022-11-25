package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mock sqlmock.Sqlmock

func TestGetEmployeeHandler(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "phone_number", "department_id", "name"}).
		AddRow("a070ba78-6ae8-11ed-8e82-64bc58925a40", "Abishek", "1234567890", "55e95991-6ae8-11ed-8e82-64bc58925a40", "Software")

	mock.ExpectQuery(selectAllEmployeeQuery).WillReturnRows(rows)

	tests := []struct {
		description string
		expCode     int
		expResp     string
	}{
		{
			description: "Case get all employee `/employee`",
			expCode:     200,
			expResp:     `[{"id":"a070ba78-6ae8-11ed-8e82-64bc58925a40","name":"Abishek","phoneNumber":"1234567890","dept":{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}}]`,
		},
	}

	for _, tc := range tests {
		mockReq, _ := http.NewRequest("GET", "/employee", nil)
		mockResp := httptest.NewRecorder()

		getEmployeeHandler(mockResp, mockReq)

		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)
		assert.Equal(t, tc.expResp, mockResp.Body.String())
	}
}

func TestGetDepartmentHandler(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("55e95991-6ae8-11ed-8e82-64bc58925a40", "Software")

	mock.ExpectQuery(selectAllDepartmentQuery).WillReturnRows(rows)

	tests := []struct {
		description string
		expCode     int
		expResp     string
	}{
		{
			description: "Case get all department `/department`",
			expCode:     200,
			expResp:     `[{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}]`,
		},
	}

	for _, tc := range tests {
		mockReq, _ := http.NewRequest("GET", "/department", nil)
		mockResp := httptest.NewRecorder()

		getDepartmentHandler(mockResp, mockReq)

		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)
		assert.Equal(t, tc.expResp, mockResp.Body.String())
	}
}

func TestPostEmployeeHandler(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	tests := []struct {
		description string
		reqBody     employee
		query       string
		args        map[string]any
		result      driver.Result
		mockErr     error
		expErr      error
		expCode     int
	}{
		{
			description: "Post a valid employee",
			reqBody:     employee{Name: "testEmp", PhoneNumber: "9876543210", Dept: department{Id: "55e95991-6ae8-11ed-8e82-64bc58925a40", Name: "testDept"}},
			query:       insertIntoEmployeeQuery,
			args:        map[string]any{"id": sqlmock.AnyArg(), "name": "testEmp", "dept_id": "55e95991-6ae8-11ed-8e82-64bc58925a40", "phoneNumber": "9876543210"},
			result:      sqlmock.NewResult(1, 1),
			mockErr:     nil,
			expErr:      nil,
			expCode:     201,
		},
		{
			description: "Post a invalid employee, phoneNumber > 10 digits",
			reqBody:     employee{Name: "testEmp", PhoneNumber: "9876543210000", Dept: department{Id: "55e95991-6ae8-11ed-8e82-64bc58925a40", Name: "testDept"}},
			query:       insertIntoEmployeeQuery,
			args:        map[string]any{"id": sqlmock.AnyArg(), "name": "testEmp", "dept_id": "55e95991-6ae8-11ed-8e82-64bc58925a40", "phoneNumber": "9876543210000"},
			result:      sqlmock.NewResult(0, 0),
			mockErr:     errors.New("sql: phoneNumber is a varchar(10) field"),
			expErr:      nil,
			expCode:     400,
		},
		{
			description: "Post a invalid employee, phoneNumber > 10 digits",
			reqBody:     employee{Name: "testEmp", PhoneNumber: "9876543210", Dept: department{Id: "-1", Name: "testDept"}},
			query:       insertIntoEmployeeQuery,
			args:        map[string]any{"id": sqlmock.AnyArg(), "name": "testEmp", "dept_id": "-1", "phoneNumber": "9876543210"},
			result:      sqlmock.NewResult(0, 0),
			mockErr:     errors.New("sql: dept_id not found"),
			expErr:      nil,
			expCode:     400,
		},
	}

	for _, tc := range tests {
		//mock http
		reqBody, err := json.Marshal(tc.reqBody)
		if err != nil {
			log.Println(err)
			return
		}

		mockReq, _ := http.NewRequest("POST", "/employee", bytes.NewReader(reqBody))
		mockResp := httptest.NewRecorder()

		mock.ExpectExec(insertIntoEmployeeQuery).WithArgs(tc.args["id"], tc.args["name"], tc.args["dept_id"], tc.args["phoneNumber"]).WillReturnResult(tc.result).WillReturnError(tc.mockErr)

		postEmployeeHandler(mockResp, mockReq)

		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Some expectations were not met %s", err)
		}
	}
}

func TestPostDepartmentHandler(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	tests := []struct {
		description string
		reqBody     department
		mockSqlFunc func(d department)
		expCode     int
	}{
		{
			description: "Post a valid department",
			reqBody:     department{Name: "NewDept"},
			mockSqlFunc: func(d department) {
				uuid, _ := uuid2.NewUUID()
				generatedUUID := uuid.String()
				mock.ExpectExec(insertIntoDepartmentQuery).WithArgs(sqlmock.AnyArg(), d.Name).
					WillReturnResult(sqlmock.NewResult(1, 1))
				rows := mock.NewRows([]string{"id", "name"}).
					AddRow(generatedUUID, d.Name)
				mock.ExpectQuery(selectDepartmentByIdQuery).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
			},
			expCode: 201,
		},
	}

	for _, tc := range tests {
		//mock http
		reqBody, err := json.Marshal(tc.reqBody)
		if err != nil {
			log.Println(err)
			return
		}

		mockReq, _ := http.NewRequest("POST", "/department", bytes.NewReader(reqBody))
		mockResp := httptest.NewRecorder()

		tc.mockSqlFunc(tc.reqBody)

		postDepartmentHandler(mockResp, mockReq)

		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Some expectations were not met %s", err)
		}
	}
}

// // Mocking request with query params
//
//	if tc.reqQueryParams != nil {
//		queryParams := mockReq.URL.Query()
//		for key, value := range tc.reqQueryParams {
//			queryParams.Set(key, value)
//		}
//		mockReq.URL.RawQuery = queryParams.Encode()
//	}
