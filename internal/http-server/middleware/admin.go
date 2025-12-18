package middleware

import (
	"context"
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/config"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

type adminContextKey string

const AdminContextKey adminContextKey = "is_admin"

// AdminMiddleware проверяет Basic Auth с логином/паролем из .env
func AdminMiddleware(adminCfg config.AdminConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Admin"`)
				response.WriteError(w, http.StatusUnauthorized, "Admin authentication required")
				return
			}

			if username != adminCfg.GetUsername() || password != adminCfg.GetPassword() {
				response.WriteError(w, http.StatusUnauthorized, "Invalid admin credentials")
				return
			}

			ctx := context.WithValue(r.Context(), AdminContextKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func IsAdmin(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(AdminContextKey).(bool)
	return ok && isAdmin
}
