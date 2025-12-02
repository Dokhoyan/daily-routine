package settings

import (
	"context"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (s *serv) GetByUserID(ctx context.Context, userID int64) (*models.UserSettings, error) {
	settings, err := s.settingsRepo.GetSettingsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}
	return settings, nil
}
