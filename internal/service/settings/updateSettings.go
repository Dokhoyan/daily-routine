package settings

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) UpdateSettings(ctx context.Context, userID int64, doNotDisturb *bool, notifyTimes *[]string) (*models.UserSettings, error) {
	currentSettings, err := s.settingsRepo.GetSettingsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current settings: %w", err)
	}

	if doNotDisturb != nil {
		currentSettings.DoNotDisturb = *doNotDisturb
	}
	if notifyTimes != nil {
		currentSettings.NotifyTimes = *notifyTimes
	}

	if err := s.settingsRepo.UpdateSettings(ctx, currentSettings); err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	return currentSettings, nil
}
