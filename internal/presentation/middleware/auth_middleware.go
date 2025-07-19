package middleware

import (
	"log/slog"
	"net/http"

	"payroll/internal/app/auth"
	"payroll/internal/app/auth/dto"
	"payroll/internal/app/common/api"
	"payroll/internal/infrastructure/log"
	"payroll/pkg/jwt"
)

type JWTAuthMiddleware func(next http.Handler) http.Handler

func WithJWTAuth(jwt jwt.JWT[*dto.JWTClaims]) JWTAuthMiddleware {
	return func(next http.Handler) http.Handler {
		var fn = func(w http.ResponseWriter, r *http.Request) {
			authString := r.Header.Get(api.HeaderAuthorization)
			if authString == "" {
				slog.WarnContext(r.Context(), "Empty Authorization header")
				api.SendUnauthenticated(w)
				return
			}

			jwtToken := authString[len("JWT "):]

			claims, err := jwt.ParseToken(jwtToken)
			if err != nil {
				slog.WarnContext(r.Context(), "Invalid JWT token", log.WithErrorAttr(err))
				api.SendUnauthenticated(w)
				return
			}

			ctx := r.Context()
			ctx = auth.WithAuthContext(ctx, claims.JWTPayloads)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
