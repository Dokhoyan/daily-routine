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

	if habit.Format != models.HabitFormatTime && habit.Format != models.HabitFormatCount && habit.Format != models.HabitFormatBinary {
		return errors.New("invalid habit format: must be 'time', 'count' or 'binary'")
	}

	if habit.Format == models.HabitFormatBinary {
		if habit.Value != 0 && habit.Value != 1 {
			return errors.New("for binary format, value must be 0 or 1")
		}
	} else if habit.Value <= 0 {
		return errors.New("value must be greater than 0")
	}

	if habit.Type != models.HabitTypeBeneficial && habit.Type != models.HabitTypeHarmful {
		return errors.New("invalid habit type: must be 'beneficial' or 'harmful'")
	}

	if err := s.habitRepo.UpdateHabit(ctx, habit); err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	return nil
}
