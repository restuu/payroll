-- name: SaveEmployeeAttendance :exec
INSERT INTO attendance (employee_id, timestamp, type, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5);