package sprint

import (
	"context"
	"fmt"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

// CheckNewHabitSprint проверяет и выполняет спринт new_habit при создании новой привычки
func (s *serv) CheckNewHabitSprint(ctx context.Context, userID int64) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// Получаем все активные спринты типа new_habit
	active := true
	sprints, err := s.sprintRepo.GetAllSprints(ctx, &active)
	if err != nil {
		return fmt.Errorf("failed to get active sprints: %w", err)
	}

	for _, sprint := range sprints {
		if sprint.Type != models.SprintTypeNewHabit {
			continue
		}

		// Получаем или создаем прогресс пользователя
		progress, err := s.sprintRepo.GetUserSprintProgress(ctx, userID, sprint.ID)
		if err != nil {
			fmt.Printf("warning: failed to load progress for sprint %d: %v\n", sprint.ID, err)
			continue
		}
		if progress == nil {
			progress = &models.UserSprintProgress{
				UserID:      userID,
				SprintID:    sprint.ID,
				CurrentDays: 0,
				IsCompleted: false,
			}
		}

		// Если уже выполнен, пропускаем
		if progress.IsCompleted {
			continue
		}

		// Выполняем спринт сразу при создании привычки
		now := time.Now()
		progress.CurrentDays = 1
		progress.IsCompleted = true
		progress.CompletedAt = &now

		// Начисляем награду
		if err := s.userRepo.AddCoins(ctx, userID, sprint.CoinsReward); err != nil {
			fmt.Printf("warning: failed to add coins to user %d: %v\n", userID, err)
		}

		// Сохраняем прогресс
		if err := s.sprintRepo.CreateOrUpdateUserSprintProgress(ctx, progress); err != nil {
			fmt.Printf("warning: failed to update progress for sprint %d: %v\n", sprint.ID, err)
		}
	}

	return nil
}

