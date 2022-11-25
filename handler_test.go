package main

import (
	"bytes"
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
		reqBody     []byte
		expCode     int
		expResp     string
	}{
		{
			description: "Case get all employee `/employee`",
			reqBody:     nil,
			expCode:     200,
			expResp:     `[{"id":"a070ba78-6ae8-11ed-8e82-64bc58925a40","name":"Abishek","phoneNumber":"1234567890","dept":{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}}]`,
		},
	}

	for _, tc := range tests {
		mockReq, _ := http.NewRequest("GET", "/employee", bytes.NewReader(tc.reqBody))
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
		reqBody     []byte
		expCode     int
		expResp     string
	}{
		{
			description: "Case get all department `/department`",
			reqBody:     nil,
			expCode:     200,
			expResp:     `[{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}]`,
		},
	}

	for _, tc := range tests {
		mockReq, _ := http.NewRequest("GET", "/department", bytes.NewReader(tc.reqBody))
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
		mockSqlFunc func(e employee)
		expCode     int
	}{
		{
			description: "Post a valid employee",
			reqBody:     employee{Name: "newGuy", PhoneNumber: "3213213211", Dept: department{Id: "55e95991-6ae8-11ed-8e82-64bc58925a40", Name: "Software"}},
			mockSqlFunc: func(e employee) {
				uuid, _ := uuid2.NewUUID()
				generatedUUID := uuid.String()
				mock.ExpectExec(insertIntoEmployeeQuery).WithArgs(sqlmock.AnyArg(), e.Name, e.Dept.Id, e.PhoneNumber).
					WillReturnResult(sqlmock.NewResult(1, 1))
				rows := mock.NewRows([]string{"id", "name", "phone_number", "department_id", "name"}).
					AddRow(generatedUUID, e.Name, e.PhoneNumber, e.Dept.Id, e.Dept.Name)
				mock.ExpectQuery(selectEmployeeByIdQuery).WithArgs(sqlmock.AnyArg()).WillReturnRows(rows)
			},
			expCode: 201,
		},
		{
			description: "Post a invalid employee, phone_number > 10 digits",
			reqBody:     employee{Name: "newGuy", PhoneNumber: "32132132324324325211", Dept: department{Id: "55e95991-6ae8-11ed-8e82-64bc58925a40", Name: "Software"}},
			mockSqlFunc: func(e employee) {
				mock.ExpectExec(insertIntoEmployeeQuery).WithArgs(sqlmock.AnyArg(), e.Name, e.Dept.Id, e.PhoneNumber).
					WillReturnError(errors.New("sql: length exceeded for varchar 10"))
			},
			expCode: 400,
		},
		{
			description: "Post a invalid employee, department id doesnt exit",
			reqBody:     employee{Name: "newGuy", PhoneNumber: "3213213211", Dept: department{Id: "-1", Name: "Software"}},
			mockSqlFunc: func(e employee) {
				mock.ExpectExec(insertIntoEmployeeQuery).WithArgs(sqlmock.AnyArg(), e.Name, e.Dept.Id, e.PhoneNumber).
					WillReturnError(errors.New("sql: foreign key reference: dept:id doesnt exits"))
			},
			expCode: 400,
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

		tc.mockSqlFunc(tc.reqBody)

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
