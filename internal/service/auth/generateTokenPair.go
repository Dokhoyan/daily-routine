package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
	"github.com/Dokhoyan/daily-routine/internal/repository/postgres"
)

func (s *serv) GenerateTokenPair(ctx context.Context, userID string, r *http.Request) (*models.TokenPair, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	activeCount, err := s.tokenRepo.GetActiveTokensCount(ctx, userIDInt)
	if err != nil {
		fmt.Printf("warning: failed to get active tokens count: %v\n", err)
	}

	maxSessions := s.authConfig.GetMaxActiveSessions()

	if activeCount >= maxSessions {
		activeTokens, err := s.tokenRepo.GetActiveTokens(ctx, userIDInt)
		if err == nil && len(activeTokens) > 0 {
			tokensToRevoke := len(activeTokens) - maxSessions + 1
			for i := 0; i < tokensToRevoke && i < len(activeTokens); i++ {
				if err := s.tokenRepo.RevokeRefreshToken(ctx, activeTokens[i].Token); err != nil {
					fmt.Printf("warning: failed to revoke old token: %v\n", err)
				}
			}
		}
	}

	tokenHash := postgres.HashToken(refreshToken)

	var deviceInfo, ipAddress *string
	if r != nil {
		device := GetDeviceInfo(r)
		ip := GetIPAddress(r)
		deviceInfo = &device
		ipAddress = &ip
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:     userIDInt,
		Token:      tokenHash,
		ExpiresAt:  time.Now().Add(refreshTokenTTL),
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
	}

	if err := s.tokenRepo.SaveRefreshToken(ctx, refreshTokenModel); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	s.logTokenAction(ctx, userIDInt, "refresh", "issued", r)
	s.logTokenAction(ctx, userIDInt, "access", "issued", r)

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
