package models

import "time"

type HabitType string

const (
	HabitTypeTime  HabitType = "time"
	HabitTypeCount HabitType = "count"
)

type Habit struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	Type         HabitType `json:"type"`
	Unit         string    `json:"unit"`
	Value        int       `json:"value"`
	IsActive     bool      `json:"is_active"`
	IsDone       bool      `json:"is_done"`
	IsBeneficial bool      `json:"is_beneficial"`
	Series       int       `json:"series"`
	CreatedAt    time.Time `json:"created_at"`
}
