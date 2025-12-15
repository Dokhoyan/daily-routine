package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/repository/postgres"
)

func (s *serv) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.Type != "refresh" {
		return "", fmt.Errorf("token is not a refresh token")
	}

	if claims.UserID == "" {
		return "", fmt.Errorf("invalid user_id in token")
	}

	userIDInt, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid user_id format: %w", err)
	}

	tokenHash := postgres.HashToken(refreshToken)
	dbToken, err := s.tokenRepo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return "", fmt.Errorf("refresh token not found or revoked: %w", err)
	}

	if dbToken.UserID != userIDInt {
		return "", fmt.Errorf("token does not belong to user")
	}

	accessToken, err := s.generateAccessToken(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}
