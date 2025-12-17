package models

import "time"

type SprintType string

const (
	SprintTypeHabitSeries   SprintType = "habit_series"   // Поддерживать серию привычки несколько дней подряд
	SprintTypeAllHabits     SprintType = "all_habits"     // Выполнять все привычки несколько дней подряд
	SprintTypeHabitIncrease SprintType = "habit_increase" // Увеличить значение любой привычки на определенное количество процентов
)

type Sprint struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Type            SprintType `json:"type"`
	TargetDays      int        `json:"target_days"`                // Количество дней для выполнения
	CoinsReward     int        `json:"coins_reward"`               // Награда в коинах
	IsActive        bool       `json:"is_active"`                  // Активен ли спринт
	HabitID         *int64     `json:"habit_id,omitempty"`         // Для habit_series и habit_increase
	MinSeries       *int       `json:"min_series,omitempty"`       // Минимальная серия для habit_series
	PercentIncrease *int       `json:"percent_increase,omitempty"` // Процент увеличения для habit_increase
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type UserSprintProgress struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	SprintID    int64      `json:"sprint_id"`
	CurrentDays int        `json:"current_days"` // Текущее количество дней выполнения
	IsCompleted bool       `json:"is_completed"` // Выполнен ли спринт
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateSprintRequest struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Type            SprintType `json:"type"`
	TargetDays      int        `json:"target_days"`
	CoinsReward     int        `json:"coins_reward"`
	HabitID         *int64     `json:"habit_id,omitempty"`
	MinSeries       *int       `json:"min_series,omitempty"`
	PercentIncrease *int       `json:"percent_increase,omitempty"`
}
