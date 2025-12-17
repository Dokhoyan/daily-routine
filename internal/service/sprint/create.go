package sprint

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) Create(ctx context.Context, req *models.CreateSprintRequest) (*models.Sprint, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Валидация в зависимости от типа спринта
	if err := s.validateSprintRequest(req); err != nil {
		return nil, fmt.Errorf("invalid sprint request: %w", err)
	}

	sprint := &models.Sprint{
		Title:           req.Title,
		Description:     req.Description,
		Type:            req.Type,
		TargetDays:      req.TargetDays,
		CoinsReward:     req.CoinsReward,
		IsActive:        true,
		HabitID:         req.HabitID,
		MinSeries:       req.MinSeries,
		PercentIncrease: req.PercentIncrease,
	}

	createdSprint, err := s.sprintRepo.CreateSprint(ctx, sprint)
	if err != nil {
		return nil, fmt.Errorf("failed to create sprint: %w", err)
	}

	return createdSprint, nil
}

func (s *serv) validateSprintRequest(req *models.CreateSprintRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.TargetDays <= 0 {
		return fmt.Errorf("target_days must be greater than 0")
	}
	if req.CoinsReward < 0 {
		return fmt.Errorf("coins_reward cannot be negative")
	}

	switch req.Type {
	case models.SprintTypeHabitSeries:
		if req.HabitID == nil {
			return fmt.Errorf("habit_id is required for habit_series sprint")
		}
		if req.MinSeries == nil || *req.MinSeries <= 0 {
			return fmt.Errorf("min_series must be greater than 0 for habit_series sprint")
		}
	case models.SprintTypeHabitIncrease:
		if req.HabitID == nil {
			return fmt.Errorf("habit_id is required for habit_increase sprint")
		}
		if req.PercentIncrease == nil || *req.PercentIncrease <= 0 {
			return fmt.Errorf("percent_increase must be greater than 0 for habit_increase sprint")
		}
	case models.SprintTypeAllHabits:
		// Нет дополнительных требований
	default:
		return fmt.Errorf("invalid sprint type: %s", req.Type)
	}

	return nil
}
