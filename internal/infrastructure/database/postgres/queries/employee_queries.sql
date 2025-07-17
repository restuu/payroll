-- name: FindEmployeeByID :one
SELECT * FROM employee
WHERE id = $1;

-- name: InsertEmployee :exec
INSERT INTO employee (username, password, salary, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5);

-- name: FindEmployeeByUsername :one
SELECT *
FROM employee
WHERE username = $1;

-- name: InsertEmployeeRole :exec
INSERT INTO employee_role(employee_id, role_id, created_by, updated_by)
VALUES (
        (select id from employee where username = sqlc.arg(username)),
        (select id from role where name = sqlc.arg(role_name)),
        sqlc.arg(created_by),
        sqlc.arg(updated_by));