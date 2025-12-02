package habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) Create(ctx context.Context, habit *models.Habit) (*models.Habit, error) {
	if habit.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	if habit.Type != models.HabitTypeTime && habit.Type != models.HabitTypeCount {
		return nil, errors.New("invalid habit type: must be 'time' or 'count'")
	}

	if habit.Value <= 0 {
		return nil, errors.New("value must be greater than 0")
	}

	if habit.IsActive == false && habit.IsDone == false {
		habit.IsActive = true
	}
	habit.IsDone = false
	habit.Series = 0

	createdHabit, err := s.habitRepo.CreateHabit(ctx, habit)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return createdHabit, nil
}
