package habit

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) ProcessDailyReset(ctx context.Context, userID int64, habits []*models.Habit) error {
	for _, habit := range habits {
		if !habit.IsActive {
			continue
		}

		if habit.IsDone {
			habit.Series++
			habit.IsDone = false
		} else {
			habit.Series = 0
		}

		if err := s.habitRepo.UpdateHabit(ctx, habit); err != nil {
			return fmt.Errorf("failed to update habit %d: %w", habit.ID, err)
		}
	}

	return nil
}
