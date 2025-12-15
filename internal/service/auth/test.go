package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

func (s *serv) GenerateTestToken(ctx context.Context, userID int64, r *http.Request) (*service.TestTokenResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		user = &models.User{
			ID:       userID,
			Username: "test_user_" + strconv.FormatInt(userID, 10),
			PhotoURL: "",
			AuthDate: time.Now(),
		}
		if err := s.userRepo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create test user: %w", err)
		}
	}

	userIDStr := strconv.FormatInt(userID, 10)
	tokenPair, err := s.GenerateTokenPair(ctx, userIDStr, r)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &service.TestTokenResponse{
		User:    user,
		Tokens:  tokenPair,
		Message: "Test tokens generated successfully",
	}, nil
}
