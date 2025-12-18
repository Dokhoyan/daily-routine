package service

import (
	"context"
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, id int64, user *models.User) error
}

type HabitService interface {
	GetByID(ctx context.Context, id int64) (*models.Habit, error)
	GetByUserID(ctx context.Context, userID int64, habitType *string, isActive *bool) ([]*models.Habit, error)
	Create(ctx context.Context, habit *models.Habit) (*models.Habit, error)
	Update(ctx context.Context, habit *models.Habit) error
	Delete(ctx context.Context, id int64) error
	ProcessDailyReset(ctx context.Context, userID int64, habits []*models.Habit) error
}

type SettingsService interface {
	GetByUserID(ctx context.Context, userID int64) (*models.UserSettings, error)
	UpdateSettings(ctx context.Context, userID int64, doNotDisturb *bool, notifyTimes *[]string) (*models.UserSettings, error)
	UpdateTimezone(ctx context.Context, userID int64, timezone string) (*models.UserSettings, error)
}

type AuthService interface {
	VerifyTelegramData(data map[string]string) bool
	GenerateTokenPair(ctx context.Context, userID string, r *http.Request) (*models.TokenPair, error)
	AuthenticateOrRegister(ctx context.Context, telegramData map[string]string, r *http.Request) (*models.AuthResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*models.UserClaims, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	RefreshTokenPair(ctx context.Context, refreshToken string, r *http.Request) (*models.TokenPair, error)
	GenerateTestToken(ctx context.Context, userID int64, r *http.Request) (*TestTokenResponse, error)
	RevokeToken(ctx context.Context, tokenString string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
}

type SprintService interface {
	GetAll(ctx context.Context, isActive *bool) ([]*models.Sprint, error)
	GetByID(ctx context.Context, id int64) (*models.Sprint, error)
	Create(ctx context.Context, req *models.CreateSprintRequest) (*models.Sprint, error)
	Update(ctx context.Context, id int64, req *models.CreateSprintRequest) (*models.Sprint, error)
	Delete(ctx context.Context, id int64) error
	GetUserProgress(ctx context.Context, userID int64) ([]*models.UserSprintProgress, error)
	CheckAndUpdateSprintProgress(ctx context.Context, userID int64) error
	CheckNewHabitSprint(ctx context.Context, userID int64) error
	ResetWeeklyProgress(ctx context.Context) error
}

type TestTokenResponse struct {
	User    *models.User      `json:"user"`
	Tokens  *models.TokenPair `json:"tokens"`
	Message string            `json:"message"`
}
