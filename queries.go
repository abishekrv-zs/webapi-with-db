package main

// For getEmployeeHandler
var selectAllEmployeeQuery = `select e.id, e.name, phone_number, department_id, d.name from employee e inner join department d on e.department_id = d.id`

// For getDepartmentHandler
var selectAllDepartmentQuery = `select id, name from department`

// For postEmployeeHandler
var insertIntoEmployeeQuery = `insert into employee values(?,?,?,?)`
var selectEmployeeByIdQuery = `select e.id, e.name, phone_number, department_id, d.name from employee e inner join department d on e.department_id = d.id where e.id=?`

// For postDepartmentHandler
var insertIntoDepartmentQuery = `insert into department values(?,?)`
var selectDepartmentByIdQuery = `select id, name from department where id=?`
