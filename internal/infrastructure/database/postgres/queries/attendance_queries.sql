-- name: SaveEmployeeAttendance :exec
INSERT INTO attendance (employee_id, timestamp, type, created_by, updated_by)
VALUES ((
			select id
			from employee
			where username = sqlc.arg(username)
			),
		sqlc.arg(timestamp),
		sqlc.arg(type),
		sqlc.arg(created_by),
		sqlc.arg(updated_by));