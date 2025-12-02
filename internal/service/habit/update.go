package habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) Update(ctx context.Context, habit *models.Habit) error {
	if habit.Title == "" {
		return errors.New("title cannot be empty")
	}

	if habit.Type != models.HabitTypeTime && habit.Type != models.HabitTypeCount {
		return errors.New("invalid habit type: must be 'time' or 'count'")
	}

	if habit.Value <= 0 {
		return errors.New("value must be greater than 0")
	}

	if err := s.habitRepo.UpdateHabit(ctx, habit); err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	return nil
}
