package repository

import (
	"context"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

// UserRepository определяет методы для доступа к данным пользователей
type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}

// HabitRepository определяет методы для доступа к данным привычек
type HabitRepository interface {
	GetHabitByID(ctx context.Context, id int64) (*models.Habit, error)
	GetHabitsByUserID(ctx context.Context, userID int64) ([]*models.Habit, error)
	CreateHabit(ctx context.Context, habit *models.Habit) (*models.Habit, error)
	UpdateHabit(ctx context.Context, habit *models.Habit) error
	DeleteHabit(ctx context.Context, id int64) error
	DeleteHabitsByUserID(ctx context.Context, userID int64) error
}

// UserSettingsRepository определяет методы для доступа к данным настроек пользователя
type UserSettingsRepository interface {
	GetSettingsByUserID(ctx context.Context, userID int64) (*models.UserSettings, error)
	CreateSettings(ctx context.Context, settings *models.UserSettings) error
	UpdateSettings(ctx context.Context, settings *models.UserSettings) error
}

// TokenRepository определяет методы для доступа к данным токенов
type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
	GetActiveTokensCount(ctx context.Context, userID int64) (int, error)
	GetActiveTokens(ctx context.Context, userID int64) ([]*models.RefreshToken, error)
	DeleteExpiredTokens(ctx context.Context) error
	LogTokenAction(ctx context.Context, log *models.TokenLog) error
	AddToBlacklist(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time, reason *string) error
	IsTokenBlacklisted(ctx context.Context, tokenHash string) (bool, error)
	DeleteExpiredBlacklistEntries(ctx context.Context) error
}
