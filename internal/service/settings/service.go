package settings

import (
	"github.com/Dokhoyan/daily-routine/internal/repository"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

type serv struct {
	settingsRepo repository.UserSettingsRepository
}

func NewService(settingsRepo repository.UserSettingsRepository) service.SettingsService {
	return &serv{
		settingsRepo: settingsRepo,
	}
}
