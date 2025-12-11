package models

import "time"

type HabitFormat string

const (
	HabitFormatTime   HabitFormat = "time"
	HabitFormatCount  HabitFormat = "count"
	HabitFormatBinary HabitFormat = "binary"
)

type HabitType string

const (
	HabitTypeBeneficial HabitType = "beneficial"
	HabitTypeHarmful    HabitType = "harmful"
)

type Habit struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	Title      string      `json:"title"`
	Format     HabitFormat `json:"format"`
	Unit       string      `json:"unit"`
	Value      int         `json:"value"`
	IsActive   bool        `json:"is_active"`
	IsDone     bool        `json:"is_done"`
	Type       HabitType   `json:"type"`
	Series     int         `json:"series"`
	CreatedAt  time.Time  `json:"created_at"`
}
