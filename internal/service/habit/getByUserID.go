package habit

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) GetByUserID(ctx context.Context, userID int64) ([]*models.Habit, error) {
	habits, err := s.habitRepo.GetHabitsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	return habits, nil
}
