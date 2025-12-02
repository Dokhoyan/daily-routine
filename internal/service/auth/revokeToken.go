package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
	"github.com/Dokhoyan/daily-routine/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

func (s *serv) RevokeToken(ctx context.Context, tokenString string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*models.UserClaims)
	if !ok {
		return fmt.Errorf("failed to extract claims")
	}

	if claims.UserID == "" {
		return fmt.Errorf("invalid user_id in token")
	}

	userIDInt, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user_id format: %w", err)
	}

	var expiresAt time.Time
	var ttl time.Duration
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
		ttl = time.Until(expiresAt)
	} else {
		if claims.Type == "access" {
			expiresAt = time.Now().Add(accessTokenTTL)
			ttl = accessTokenTTL
		} else {
			expiresAt = time.Now().Add(refreshTokenTTL)
			ttl = refreshTokenTTL
		}
	}

	tokenHash := postgres.HashToken(tokenString)

	if claims.Type == "refresh" {
		if err := s.tokenRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
			fmt.Printf("warning: failed to revoke refresh token: %v\n", err)
		}
	}

	reason := "user_revoked"
	if err := s.tokenRepo.AddToBlacklist(ctx, tokenHash, userIDInt, expiresAt, &reason); err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	if s.tokenCache != nil && ttl > 0 {
		s.tokenCache.AddToBlacklist(ctx, tokenHash, ttl)
	}

	s.logTokenAction(ctx, userIDInt, claims.Type, "revoked", nil)

	return nil
}

func (s *serv) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	if ctx == nil {
		ctx = context.Background()
	}

	activeTokens, err := s.tokenRepo.GetActiveTokens(ctx, userID)
	if err != nil {
		fmt.Printf("warning: failed to get active tokens: %v\n", err)
	}

	if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}

	if s.tokenCache != nil {
		now := time.Now()
		for _, token := range activeTokens {
			if token.ExpiresAt.After(now) {
				ttl := time.Until(token.ExpiresAt)
				if ttl > 0 {
					s.tokenCache.AddToBlacklist(ctx, token.Token, ttl)
				}
			}
		}
	}

	s.logTokenAction(ctx, userID, "all", "revoked", nil)

	return nil
}
