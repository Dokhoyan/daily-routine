package habit

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) GetByID(ctx context.Context, id int64) (*models.Habit, error) {
	habit, err := s.habitRepo.GetHabitByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}
	return habit, nil
}
