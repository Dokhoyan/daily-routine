package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) Update(ctx context.Context, id int64, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	currentUser, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Username != "" {
		currentUser.Username = user.Username
	}
	if user.FirstName != "" {
		currentUser.FirstName = user.FirstName
	}
	if user.PhotoURL != "" {
		currentUser.PhotoURL = user.PhotoURL
	}

	currentUser.ID = id
	if err := s.userRepo.UpdateUser(ctx, currentUser); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
