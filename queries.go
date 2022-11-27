package main

var getAllEmployeeQuery = `select e.id, e.name, phone_number, department_id, d.name from employee e inner join department d on department_id = d.id`

var getAllEmployeeByDeptIdQuery = getAllEmployeeQuery + ` where department_id=?`

var getEmployeeByIdQuery = getAllEmployeeQuery + ` where e.id=?`
