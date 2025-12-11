package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) AuthenticateOrRegister(ctx context.Context, telegramData map[string]string, r *http.Request) (*models.AuthResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if !s.VerifyTelegramData(telegramData) {
		return nil, fmt.Errorf("telegram data verification failed")
	}

	var userID int64
	if _, err := fmt.Sscanf(telegramData["id"], "%d", &userID); err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	var authDate time.Time
	if authDateStr, ok := telegramData["auth_date"]; ok && authDateStr != "" {
		var timestamp int64
		if _, err := fmt.Sscanf(authDateStr, "%d", &timestamp); err == nil {
			authDate = time.Unix(timestamp, 0)
		} else {
			authDate = time.Now()
		}
	} else {
		authDate = time.Now()
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if !strings.Contains(err.Error(), "user not found") {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		user = &models.User{
			ID:        userID,
			Username:  telegramData["username"],
			FirstName: telegramData["first_name"],
			PhotoURL:  telegramData["photo_url"],
			AuthDate:  authDate,
			TokenTG:   "",
		}

		if err := s.userRepo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		needsUpdate := false
		if user.Username == "" && telegramData["username"] != "" {
			user.Username = telegramData["username"]
			needsUpdate = true
		}
		if user.FirstName == "" && telegramData["first_name"] != "" {
			user.FirstName = telegramData["first_name"]
			needsUpdate = true
		}
		if user.PhotoURL == "" && telegramData["photo_url"] != "" {
			user.PhotoURL = telegramData["photo_url"]
			needsUpdate = true
		}
		if needsUpdate {
			user.AuthDate = authDate
			if err := s.userRepo.UpdateUser(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
		}
	}

	tokenPair, err := s.GenerateTokenPair(ctx, telegramData["id"], r)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		User:      user,
		TokenPair: tokenPair,
	}, nil
}
