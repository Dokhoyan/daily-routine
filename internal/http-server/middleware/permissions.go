package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/repository"
)

func RequireOwnership(userRepo repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireUserID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDFromToken, err := GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userIDFromRequest := r.URL.Query().Get("user_id")
			if userIDFromRequest == "" {
				next.ServeHTTP(w, r)
				return
			}

			if userIDFromToken != userIDFromRequest {
				http.Error(w, "Forbidden: access denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireUserIDFromPath() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDFromTokenStr, err := GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userIDFromToken, err := strconv.ParseInt(userIDFromTokenStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "token_user_id", userIDFromToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromTokenContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value("token_user_id").(int64)
	if !ok {
		return 0, fmt.Errorf("token_user_id not found in context")
	}
	return userID, nil
}
