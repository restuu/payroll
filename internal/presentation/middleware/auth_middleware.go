package middleware

import (
	"log/slog"
	"net/http"
	"slices"

	"payroll/internal/app/auth"
	"payroll/internal/app/auth/dto"
	"payroll/internal/app/common/api"
	"payroll/internal/infrastructure/database/postgres/repository"
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

type IsAdminMiddleware func(next http.Handler) http.Handler

func IsAdmin(authRepo auth.AuthRepository) IsAdminMiddleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			acct := auth.GetAuthContext(ctx)
			if acct.Username == "" {
				api.SendUnauthenticated(w)
				return
			}

			roles, err := authRepo.FindRolesByUsername(ctx, acct.Username)
			if err != nil {
				slog.ErrorContext(ctx, "IsAdmin middleware failed find roles", log.WithErrorAttr(err))
				api.SendInternalServerError(w)
				return
			}

			if !slices.Contains(roles, repository.RoleTypeADMIN) {
				slog.WarnContext(ctx, "IsAdmin invalid role", slog.Any("role", roles))
				api.SendUnauthenticated(w)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
