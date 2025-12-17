package sprint

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) GetByID(ctx context.Context, id int64) (*models.Sprint, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sprint, err := s.sprintRepo.GetSprintByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprint: %w", err)
	}

	return sprint, nil
}

