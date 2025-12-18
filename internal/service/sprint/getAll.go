package sprint

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) GetAll(ctx context.Context, isActive *bool) ([]*models.Sprint, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sprints, err := s.sprintRepo.GetAllSprints(ctx, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprints: %w", err)
	}

	return sprints, nil
}


