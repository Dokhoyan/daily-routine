package sprint

import (
	"context"
	"fmt"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) GetUserProgress(ctx context.Context, userID int64) ([]*models.UserSprintProgress, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	progresses, err := s.sprintRepo.GetUserSprintProgresses(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	return progresses, nil
}

func (s *serv) CheckAndUpdateSprintProgress(ctx context.Context, userID int64) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// Получаем все активные спринты
	active := true
	sprints, err := s.sprintRepo.GetAllSprints(ctx, &active)
	if err != nil {
		return fmt.Errorf("failed to get active sprints: %w", err)
	}

	for _, sprint := range sprints {
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

		// Проверяем выполнение спринта в зависимости от типа
		completed, err := s.checkSprintCompletion(ctx, userID, sprint, progress)
		if err != nil {
			fmt.Printf("warning: failed to check sprint %d completion: %v\n", sprint.ID, err)
			continue
		}

		if completed {
			// Увеличиваем счетчик дней только если еще не достигнута цель
			if !progress.IsCompleted {
				progress.CurrentDays++

				// Проверяем, достигнута ли цель
				if progress.CurrentDays >= sprint.TargetDays {
					now := time.Now()
					progress.IsCompleted = true
					progress.CompletedAt = &now

					// Начисляем награду только один раз
					if err := s.userRepo.AddCoins(ctx, userID, sprint.CoinsReward); err != nil {
						fmt.Printf("warning: failed to add coins to user %d: %v\n", userID, err)
					}
				}

				// Сохраняем прогресс
				if err := s.sprintRepo.CreateOrUpdateUserSprintProgress(ctx, progress); err != nil {
					fmt.Printf("warning: failed to update progress for sprint %d: %v\n", sprint.ID, err)
				}
			}
		}
	}

	return nil
}

func (s *serv) checkSprintCompletion(ctx context.Context, userID int64, sprint *models.Sprint, progress *models.UserSprintProgress) (bool, error) {
	switch sprint.Type {
	case models.SprintTypeHabitSeries:
		return s.checkHabitSeriesSprint(ctx, userID, sprint, progress)
	case models.SprintTypeAllHabits:
		return s.checkAllHabitsSprint(ctx, userID, sprint, progress)
	case models.SprintTypeHabitIncrease:
		return s.checkHabitIncreaseSprint(ctx, userID, sprint, progress)
	default:
		return false, fmt.Errorf("unknown sprint type: %s", sprint.Type)
	}
}

// checkHabitSeriesSprint проверяет, поддерживает ли пользователь серию привычки
func (s *serv) checkHabitSeriesSprint(ctx context.Context, userID int64, sprint *models.Sprint, progress *models.UserSprintProgress) (bool, error) {
	if sprint.HabitID == nil || sprint.MinSeries == nil {
		return false, fmt.Errorf("habit_id and min_series are required for habit_series sprint")
	}

	habit, err := s.habitRepo.GetHabitByID(ctx, *sprint.HabitID)
	if err != nil {
		return false, fmt.Errorf("failed to get habit: %w", err)
	}

	// Проверяем, что привычка принадлежит пользователю
	if habit.UserID != userID {
		return false, fmt.Errorf("habit does not belong to user")
	}

	// Проверяем, что серия привычки >= min_series
	return habit.Series >= *sprint.MinSeries, nil
}

// checkAllHabitsSprint проверяет, выполнены ли все активные привычки пользователя
func (s *serv) checkAllHabitsSprint(ctx context.Context, userID int64, sprint *models.Sprint, progress *models.UserSprintProgress) (bool, error) {
	active := true
	habits, err := s.habitRepo.GetHabitsByUserID(ctx, userID, nil, &active)
	if err != nil {
		return false, fmt.Errorf("failed to get habits: %w", err)
	}

	if len(habits) == 0 {
		return false, nil
	}

	// Проверяем, что все активные привычки выполнены
	for _, habit := range habits {
		if !habit.IsDone {
			return false, nil
		}
	}

	return true, nil
}

// checkHabitIncreaseSprint проверяет, увеличилось ли значение привычки на нужный процент
func (s *serv) checkHabitIncreaseSprint(ctx context.Context, userID int64, sprint *models.Sprint, progress *models.UserSprintProgress) (bool, error) {
	if sprint.HabitID == nil || sprint.PercentIncrease == nil {
		return false, fmt.Errorf("habit_id and percent_increase are required for habit_increase sprint")
	}

	habit, err := s.habitRepo.GetHabitByID(ctx, *sprint.HabitID)
	if err != nil {
		return false, fmt.Errorf("failed to get habit: %w", err)
	}

	// Проверяем, что привычка принадлежит пользователю
	if habit.UserID != userID {
		return false, fmt.Errorf("habit does not belong to user")
	}

	// Вычисляем требуемое увеличение
	// Если это первый день, сохраняем базовое значение
	// Для упрощения считаем, что значение должно увеличиться на X% от исходного значения
	// В реальности нужно хранить базовое значение при старте спринта

	// Простая проверка: если current_value > 0 и value увеличилось
	// Для более точной проверки нужно хранить начальное значение в progress или отдельной таблице
	// Пока используем упрощенную логику: проверяем, что current_value достиг определенного процента от value
	if habit.Value > 0 {
		requiredIncrease := habit.Value * *sprint.PercentIncrease / 100
		return habit.CurrentValue >= requiredIncrease, nil
	}

	return false, nil
}
