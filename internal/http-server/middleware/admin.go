package middleware

import (
	"net/http"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

func AdminMiddleware(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDStr, err := GetUserIDFromContext(r.Context())
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
				return
			}

			user, err := userService.GetByID(r.Context(), userID)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "User not found")
				return
			}

			if !user.IsAdmin {
				response.WriteError(w, http.StatusForbidden, "Admin access required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
