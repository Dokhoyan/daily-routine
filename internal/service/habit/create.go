package habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/logger"
	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) Create(ctx context.Context, habit *models.Habit) (*models.Habit, error) {
	logger.Infof("habit service create: started for userID=%d, title=%s", habit.UserID, habit.Title)

	if habit.Title == "" {
		logger.Warn("habit service create: validation failed - title is empty")
		return nil, errors.New("title cannot be empty")
	}

	if habit.Format != models.HabitFormatTime && habit.Format != models.HabitFormatCount && habit.Format != models.HabitFormatBinary {
		logger.Warnf("habit service create: validation failed - invalid format: %s", habit.Format)
		return nil, errors.New("invalid habit format: must be 'time', 'count' or 'binary'")
	}

	if habit.Format == models.HabitFormatBinary {
		if habit.Value != 0 && habit.Value != 1 {
			logger.Warnf("habit service create: validation failed - binary format value must be 0 or 1, got: %d", habit.Value)
			return nil, errors.New("for binary format, value must be 0 or 1")
		}
	} else if habit.Value <= 0 {
		logger.Warnf("habit service create: validation failed - value must be > 0, got: %d", habit.Value)
		return nil, errors.New("value must be greater than 0")
	}

	if habit.CurrentValue < 0 {
		logger.Warnf("habit service create: validation failed - current_value < 0: %d", habit.CurrentValue)
		return nil, errors.New("current_value cannot be negative")
	}

	if habit.CurrentValue > habit.Value {
		logger.Warnf("habit service create: validation failed - current_value > value: %d > %d", habit.CurrentValue, habit.Value)
		return nil, errors.New("current_value cannot be greater than value")
	}

	if habit.Type != models.HabitTypeBeneficial && habit.Type != models.HabitTypeHarmful {
		logger.Warnf("habit service create: validation failed - invalid type: %s", habit.Type)
		return nil, errors.New("invalid habit type: must be 'beneficial' or 'harmful'")
	}

	logger.Info("habit service create: validation passed")

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

	logger.Infof("habit service create: calling repository CreateHabit")
	createdHabit, err := s.habitRepo.CreateHabit(ctx, habit)
	if err != nil {
		logger.Errorf("habit service create: repository error: %v", err)
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	logger.Infof("habit service create: habit created with id=%d", createdHabit.ID)
	return createdHabit, nil
}
