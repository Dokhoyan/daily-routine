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

	if habit.Type != models.HabitTypeTime && habit.Type != models.HabitTypeCount && habit.Type != models.HabitTypeBinary {
		return nil, errors.New("invalid habit type: must be 'time', 'count' or 'binary'")
	}

	if habit.Type == models.HabitTypeBinary {
		if habit.Value != 0 && habit.Value != 1 {
			return nil, errors.New("for binary type, value must be 0 or 1")
		}
	} else if habit.Value <= 0 {
		return nil, errors.New("value must be greater than 0")
	}

	if !habit.IsActive && !habit.IsDone {
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
