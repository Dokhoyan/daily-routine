package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/models"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
	ClaimsKey ContextKey = "claims"
)

func AuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims, err := authService.ValidateToken(r.Context(), tokenString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
				return
			}

			if claims.Type != "access" {
				http.Error(w, "Token is not an access token", http.StatusUnauthorized)
				return
			}

			if claims.UserID == "" {
				http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ClaimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}

func GetClaimsFromContext(ctx context.Context) (*models.UserClaims, error) {
	claims, ok := ctx.Value(ClaimsKey).(*models.UserClaims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

func OptionalAuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString := parts[1]

					claims, err := authService.ValidateToken(r.Context(), tokenString)
					if err == nil {
						if claims.Type == "access" && claims.UserID != "" {
							ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
							ctx = context.WithValue(ctx, ClaimsKey, claims)
							r = r.WithContext(ctx)
						}
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
