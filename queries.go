package main

// For getEmployeeHandler
var selectAllEmployeeQuery = `select e.id, e.name, phone_number, department_id, d.name from employee e inner join department d on e.department_id = d.id`

// For getDepartmentHandler
var selectAllDepartmentQuery = `select id, name from department`
