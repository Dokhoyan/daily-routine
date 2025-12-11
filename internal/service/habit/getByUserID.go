package habit

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) GetByUserID(ctx context.Context, userID int64, habitType *string, isActive *bool) ([]*models.Habit, error) {
	var typeFilter *models.HabitType
	if habitType != nil {
		ht := models.HabitType(*habitType)
		typeFilter = &ht
	}

	habits, err := s.habitRepo.GetHabitsByUserID(ctx, userID, typeFilter, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	return habits, nil
}
