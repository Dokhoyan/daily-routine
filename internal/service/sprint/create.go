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
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		TargetDays:  req.TargetDays,
		CoinsReward: req.CoinsReward,
		IsActive:    true,
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
	// target_days проверяется отдельно для all_habits
	if req.CoinsReward < 0 {
		return fmt.Errorf("coins_reward cannot be negative")
	}

	switch req.Type {
	case models.SprintTypeAllHabits:
		if req.TargetDays <= 0 {
			return fmt.Errorf("target_days must be greater than 0 for all_habits sprint")
		}
	case models.SprintTypeNewHabit:
		// Нет дополнительных требований, target_days не нужен
	default:
		return fmt.Errorf("invalid sprint type: %s", req.Type)
	}

	return nil
}
