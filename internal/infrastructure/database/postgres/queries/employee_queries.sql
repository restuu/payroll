-- name: FindEmployeeByID :one
SELECT * FROM employee
WHERE id = $1;

-- name: InsertEmployee :exec
INSERT INTO employee (username, password, salary, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5);

