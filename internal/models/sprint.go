package models

import "time"

type SprintType string

const (
	SprintTypeAllHabits  SprintType = "all_habits"  // Выполнять все привычки несколько дней подряд
	SprintTypeNewHabit   SprintType = "new_habit"   // Создать новую привычку
)

type Sprint struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        SprintType `json:"type"`
	TargetDays  int        `json:"target_days"`  // Количество дней для выполнения (для all_habits)
	CoinsReward int        `json:"coins_reward"` // Награда в коинах
	IsActive    bool       `json:"is_active"`    // Активен ли спринт
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        SprintType `json:"type"`
	TargetDays  int        `json:"target_days"` // Для all_habits
	CoinsReward int        `json:"coins_reward"`
}
