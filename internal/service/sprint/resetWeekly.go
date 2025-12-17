package sprint

import (
	"context"
	"fmt"
)

func (s *serv) ResetWeeklyProgress(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// Обнуляем прогресс для всех спринтов (еженедельный сброс запускается по крону по понедельникам)
	if err := s.sprintRepo.ResetAllUserSprintProgresses(ctx); err != nil {
		return fmt.Errorf("failed to reset weekly progress: %w", err)
	}

	return nil
}
