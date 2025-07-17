package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"payroll/internal/app/common/apperror"
)

type HandleHTTPFunc[T any] func(*http.Request) (T, error)

func HandleHTTP[T any](fn HandleHTTPFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data T
			err  error
		)

		data, err = fn(r)
		if err != nil {
			var e *apperror.AppError
			if errors.As(err, &e) {
				httpCode := e.HttpCode()
				if httpCode == 0 {
					httpCode = http.StatusInternalServerError
				}

				http.Error(w, e.Error(), httpCode)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]any{
			"code": "SUCCESS",
			"data": data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
