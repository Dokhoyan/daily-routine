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

	if habit.Format != models.HabitFormatTime && habit.Format != models.HabitFormatCount && habit.Format != models.HabitFormatBinary {
		return nil, errors.New("invalid habit format: must be 'time', 'count' or 'binary'")
	}

	if habit.Format == models.HabitFormatBinary {
		if habit.Value != 0 && habit.Value != 1 {
			return nil, errors.New("for binary format, value must be 0 or 1")
		}
	} else if habit.Value <= 0 {
		return nil, errors.New("value must be greater than 0")
	}

	if habit.CurrentValue < 0 {
		return nil, errors.New("current_value cannot be negative")
	}

	if habit.CurrentValue > habit.Value {
		return nil, errors.New("current_value cannot be greater than value")
	}

	if habit.Type != models.HabitTypeBeneficial && habit.Type != models.HabitTypeHarmful {
		return nil, errors.New("invalid habit type: must be 'beneficial' or 'harmful'")
	}

	if !habit.IsActive && !habit.IsDone {
		habit.IsActive = true
	}

	if habit.Type == models.HabitTypeHarmful {
		habit.IsDone = true
		habit.Series = 0
	} else {
		habit.IsDone = false
		habit.Series = 0
	}

	createdHabit, err := s.habitRepo.CreateHabit(ctx, habit)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	return createdHabit, nil
}
