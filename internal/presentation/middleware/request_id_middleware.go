package middleware

import (
	"context"
	"net/http"

	"payroll/internal/infrastructure/log"

	"github.com/google/uuid"
)

func RequestID() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			ctx := context.WithValue(r.Context(), log.RequestIDCtxKey, requestID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func GetRequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(log.RequestIDCtxKey).(string)
	return id
}
