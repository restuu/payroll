package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

//go:generate go tool enumer -type=Code -json -transform=kebab -trimprefix Code
type Code int

const (
	_ Code = iota
	BadRequest
	Unauthorized
	Forbidden
	NotFound
	Conflict
	InternalServerError Code = 9999
)

type AppError struct {
	code     Code
	httpCode int
	err      error
	msg      string
}

func (a *AppError) Code() Code {
	return a.code
}

func (a *AppError) HttpCode() int {
	return a.httpCode
}

func (a *AppError) Err() error {
	return a.err
}

func (a *AppError) Msg() string {
	return a.msg
}

func (a *AppError) WithError(err error) *AppError {
	x := *a
	x.err = errors.Join(x.err, err)

	return &x
}

func (a *AppError) Error() string {
	return fmt.Sprintf("code: %d, http code: %d, message: %s, error: %v", a.code, a.httpCode, a.msg, a.err)
}

var (
	ErrConflict       = &AppError{code: Conflict, msg: "Conflict", httpCode: http.StatusConflict}
	ErrNotFound       = &AppError{code: NotFound, msg: "Not Found", httpCode: http.StatusNotFound}
	ErrUnauthorized   = &AppError{code: Unauthorized, msg: "Unauthorized", httpCode: http.StatusUnauthorized}
	ErrInternalServer = &AppError{code: InternalServerError, msg: "Internal Server Error", httpCode: http.StatusInternalServerError}
)

func WrapWith(appError *AppError, err error) *AppError {
	appError.err = errors.Join(appError.err, err)

	return appError
}
