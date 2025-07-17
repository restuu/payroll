-- +goose Up
-- +goose StatementBegin
-- Create a new ENUM type for attendance events. This must be done
-- outside of a transaction.
CREATE TYPE attendance_type AS ENUM ('CLOCK_IN', 'CLOCK_OUT', 'OVERTIME_IN', 'OVERTIME_OUT');
-- Create a new ENUM type for user roles.
CREATE TYPE role_type AS ENUM ('SUPERUSER', 'ADMIN', 'USER');
-- +goose StatementEnd

-- +goose StatementBegin
-- Create a function to automatically set updated_at on table updates.
-- This function will be used by triggers on each table.
CREATE OR REPLACE FUNCTION set_updated_at_now()
	RETURNS TRIGGER AS
$$
BEGIN
	-- If the updated_at column is not being changed in the UPDATE statement,
	-- set it to the current timestamp. This allows for manual overrides.
	IF NEW.updated_at = OLD.updated_at THEN
		NEW.updated_at = now();
	END IF;
	RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employee
(
	id         serial         not null primary key,
	username   varchar(20)    not null unique,
	password   text           not null,
	salary     numeric(13, 2) not null,
	created_at timestamptz    not null default current_timestamp,
	created_by text           not null,
	updated_at timestamptz    not null default current_timestamp,
	updated_by text           not null
);

CREATE TRIGGER set_updated_at_employee
	BEFORE UPDATE
	ON employee
	FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_now();

CREATE TABLE IF NOT EXISTS role
(
    id         serial      not null primary key,
    name       role_type   not null unique,
    created_at timestamptz not null default current_timestamp,
    created_by text        not null,
    updated_at timestamptz not null default current_timestamp,
    updated_by text        not null
);

CREATE TRIGGER set_updated_at_role
    BEFORE UPDATE
    ON role
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_now();

INSERT INTO role(name, created_by, updated_by)
VALUES
	('SUPERUSER', 'SYSTEM', 'SYSTEM'),
	('ADMIN', 'SYSTEM', 'SYSTEM'),
	('USER', 'SYSTEM', 'SYSTEM');

CREATE TABLE IF NOT EXISTS employee_role
(
    employee_id int         not null references employee (id) on delete cascade,
    role_id     int         not null references role (id) on delete cascade,
    created_at  timestamptz not null default current_timestamp,
    created_by  text        not null,
    updated_at  timestamptz not null default current_timestamp,
    updated_by  text        not null,
    PRIMARY KEY (employee_id, role_id)
);

CREATE TRIGGER set_updated_at_employee_role
	BEFORE UPDATE
	ON employee_role
	FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_now();

CREATE TABLE IF NOT EXISTS attendance
(
	"id"          bigserial       not null primary key,
	"employee_id" int             not null references employee (id),
	"timestamp"   timestamp       not null,
	"type"        attendance_type not null,
	"date"        date GENERATED ALWAYS AS ("timestamp"::date) STORED,
	created_at    timestamptz     not null default current_timestamp,
	created_by    text            not null,
	updated_at    timestamptz     not null default current_timestamp,
	updated_by    text            not null,

	-- An employee can only have one of each attendance type per day.
	CONSTRAINT uq_attendance_employee_date_type UNIQUE (employee_id, "date", "type")
);

CREATE TRIGGER set_updated_at_attendance
	BEFORE UPDATE
	ON attendance
	FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop tables in the reverse order of creation to respect foreign key constraints.
-- Triggers are dropped automatically when the table is dropped.
DROP TABLE IF EXISTS attendance;
DROP TABLE IF EXISTS employee_role;
DROP TABLE IF EXISTS role;
DROP TABLE IF EXISTS employee;
-- +goose StatementEnd
-- +goose StatementBegin
-- Drop the trigger function
DROP FUNCTION IF EXISTS set_updated_at_now();
-- +goose StatementEnd
-- +goose StatementBegin
-- Drop the custom ENUM type. This must also be done outside a transaction.
DROP TYPE IF EXISTS attendance_type;
DROP TYPE IF EXISTS role_type;
-- +goose StatementEnd
