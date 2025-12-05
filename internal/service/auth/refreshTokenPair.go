package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/models"
	"github.com/Dokhoyan/daily-routine/internal/repository/postgres"
)

func (s *serv) RefreshTokenPair(ctx context.Context, refreshToken string, r *http.Request) (*models.TokenPair, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.Type != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	if claims.UserID == "" {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	userIDInt, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	tokenHash := postgres.HashToken(refreshToken)
	dbToken, err := s.tokenRepo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or revoked: %w", err)
	}

	if dbToken.UserID != userIDInt {
		return nil, fmt.Errorf("token does not belong to user")
	}

	tokenPair, err := s.generateTokenPairWithSessionCheck(ctx, claims.UserID, r, false)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	if err := s.tokenRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {

		fmt.Printf("warning: failed to revoke old refresh token: %v\n", err)
	} else {
		s.logTokenAction(ctx, userIDInt, "refresh", "revoked", r)
	}

	s.logTokenAction(ctx, userIDInt, "refresh", "refreshed", r)

	return tokenPair, nil
}
