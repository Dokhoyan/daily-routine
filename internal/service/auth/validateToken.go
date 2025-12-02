package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
	"github.com/Dokhoyan/daily-routine/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

func (s *serv) ValidateToken(ctx context.Context, tokenString string) (*models.UserClaims, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	tokenHash := postgres.HashToken(tokenString)

	if s.tokenCache != nil && s.tokenCache.IsBlacklisted(ctx, tokenHash) {
		return nil, fmt.Errorf("token has been revoked")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(*models.UserClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	isBlacklisted, err := s.tokenRepo.IsTokenBlacklisted(ctx, tokenHash)
	if err != nil {
		fmt.Printf("warning: failed to check blacklist: %v\n", err)
	} else if isBlacklisted {
		if s.tokenCache != nil && claims.ExpiresAt != nil {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				s.tokenCache.AddToBlacklist(ctx, tokenHash, ttl)
			}
		}
		return nil, fmt.Errorf("token has been revoked")
	}

	return claims, nil
}
