package sprint

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) Update(ctx context.Context, id int64, req *models.CreateSprintRequest) (*models.Sprint, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Проверяем существование спринта
	existingSprint, err := s.sprintRepo.GetSprintByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}

	// Валидация
	if err := s.validateSprintRequest(req); err != nil {
		return nil, fmt.Errorf("invalid sprint request: %w", err)
	}

	// Обновляем поля
	existingSprint.Title = req.Title
	existingSprint.Description = req.Description
	existingSprint.Type = req.Type
	existingSprint.TargetDays = req.TargetDays
	existingSprint.CoinsReward = req.CoinsReward

	if err := s.sprintRepo.UpdateSprint(ctx, existingSprint); err != nil {
		return nil, fmt.Errorf("failed to update sprint: %w", err)
	}

	return existingSprint, nil
}


