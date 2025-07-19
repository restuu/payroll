package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"payroll/internal/app/common/apperror"
)

type HandleHTTPFunc[T any] func(*http.Request) (T, error)

type HTTPResponse struct {
	Code  string `json:"code"`
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

func HandleHTTP[T any](fn HandleHTTPFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data     T
			err      error
			response HTTPResponse
		)

		data, err = fn(r)
		if err != nil {
			var (
				httpCode int
				errMsg   string
			)

			var e *apperror.AppError
			if errors.As(err, &e) {
				httpCode = e.HttpCode()
				if httpCode == 0 {
					httpCode = http.StatusInternalServerError
				}
				errMsg = e.Msg()
			} else {
				httpCode = http.StatusInternalServerError
				errMsg = "Internal Server Error"
			}

			response.Code = "FAILED"
			response.Error = errMsg

			w.Header().Set(HeaderContentType, MimeApplicationJSON)
			w.WriteHeader(httpCode)
			_ = json.NewEncoder(w).Encode(response)

			return
		}

		response.Code = "SUCCESS"
		response.Data = data

		w.Header().Set(HeaderContentType, MimeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}

func SendUnauthenticated(w http.ResponseWriter) {
	w.Header().Set(HeaderContentType, MimeApplicationJSON)
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(HTTPResponse{
		Code:  "FAILED",
		Error: "Unauthenticated",
	})
}
