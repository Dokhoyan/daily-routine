package habit

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	// Проверяем спринт new_habit
	s.checkNewHabitSprint(ctx, habit.UserID)

	return createdHabit, nil
}

func (s *serv) checkNewHabitSprint(ctx context.Context, userID int64) {
	active := true
	sprints, err := s.sprintRepo.GetAllSprints(ctx, &active)
	if err != nil {
		logger.Warnf("habit service: failed to get sprints: %v", err)
		return
	}

	for _, sprint := range sprints {
		if sprint.Type != models.SprintTypeNewHabit {
			continue
		}

		progress, err := s.sprintRepo.GetUserSprintProgress(ctx, userID, sprint.ID)
		if err != nil {
			logger.Warnf("habit service: failed to get sprint progress: %v", err)
			continue
		}

		if progress != nil && progress.IsCompleted {
			continue
		}

		// Создаём или обновляем прогресс как выполненный
		now := time.Now()
		newProgress := &models.UserSprintProgress{
			UserID:      userID,
			SprintID:    sprint.ID,
			CurrentDays: 1,
			IsCompleted: true,
			CompletedAt: &now,
		}
		if progress != nil {
			newProgress.ID = progress.ID
		}

		if err := s.sprintRepo.CreateOrUpdateUserSprintProgress(ctx, newProgress); err != nil {
			logger.Warnf("habit service: failed to update sprint progress: %v", err)
			continue
		}

		// Начисляем награду
		if err := s.userRepo.AddCoins(ctx, userID, sprint.CoinsReward); err != nil {
			logger.Warnf("habit service: failed to add coins: %v", err)
		} else {
			logger.Infof("habit service: new_habit sprint completed for user %d, awarded %d coins", userID, sprint.CoinsReward)
		}
	}
}
