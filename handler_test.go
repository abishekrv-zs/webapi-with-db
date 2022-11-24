package main

import (
	"bytes"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

var mock sqlmock.Sqlmock

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

//
//func TestPostEmployeeHandler(t *testing.T) {
//	db, err = getDbConnection("mysql", "abishek:mypassword^123@tcp(127.0.0.1:3306)/sample_db")
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	defer db.Close()
//
//	tests := []struct {
//		description string
//		reqBody     []byte
//		expCode     int
//		expResp     string
//	}{
//		{
//			description: "Post invalid employee(phonenumber >10 digits)",
//			reqBody:     []byte(`{"name":"Ram","phoneNumber":"987652132143210","dept":{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}}`),
//			expCode:     400,
//			expResp:     `{"error":"<error message>"}`,
//		},
//		{
//			description: "Post invalid employee(department doesn't exist)",
//			reqBody:     []byte(`{"name": "Ram","phoneNumber": "1234567111","dept": {"id": "-1"}}`),
//			expCode:     400,
//			expResp:     `{"error":"<error message>"}`,
//		},
//		{
//			description: "Post invalid employee(name is integer)",
//			reqBody:     []byte(`{"name":22,"phoneNumber":"1231231231","dept":{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}}`),
//			expCode:     400,
//			expResp:     `{"error":"<error message>}`,
//		},
//		{
//			description: "Post a valid employee",
//			reqBody:     []byte(`{"name":"newGuy","phoneNumber":"3213213211","dept":{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}}`),
//			expCode:     201,
//			expResp:     `{"id":"f8458145-548c-4e6b-98ea-62007ee779d6","name":"newGuy","phoneNumber":"3213213211","dept":{"id":"55e95991-6ae8-11ed-8e82-64bc58925a40","name":"Software"}}`,
//		},
//	}
//
//	for _, tc := range tests {
//		mockReq, _ := http.NewRequest("POST", "/employee", bytes.NewReader(tc.reqBody))
//		mockResp := httptest.NewRecorder()
//
//		postEmployeeHandler(mockResp, mockReq)
//
//		assert.Equal(t, tc.expCode, mockResp.Code, tc.description)
//		if tc.expCode >= 200 && tc.expCode < 300 { // If request is valid
//			assert.Equal(t, tc.expResp, mockResp.Body.String())
//		}
//	}
//}

//// Mocking request with query params
//if tc.reqQueryParams != nil {
//	queryParams := mockReq.URL.Query()
//	for key, value := range tc.reqQueryParams {
//		queryParams.Set(key, value)
//	}
//	mockReq.URL.RawQuery = queryParams.Encode()
//}
