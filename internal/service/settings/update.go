package settings

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
)

func validateTimezone(tz string) (string, error) {
	if strings.TrimSpace(tz) == "" {
		return "", fmt.Errorf("timezone cannot be empty")
	}

	tz = strings.TrimSpace(tz)

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %w", err)
	}

	return loc.String(), nil
}

func (s *serv) UpdateTimezone(ctx context.Context, userID int64, timezone string) (*models.UserSettings, error) {
	normalizedTZ, err := validateTimezone(timezone)
	if err != nil {
		return nil, fmt.Errorf("timezone validation failed: %w", err)
	}

	currentSettings, err := s.settingsRepo.GetSettingsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current settings: %w", err)
	}

	currentSettings.Timezone = normalizedTZ

	if err := s.settingsRepo.UpdateSettings(ctx, currentSettings); err != nil {
		return nil, fmt.Errorf("failed to update timezone: %w", err)
	}

	return currentSettings, nil
}
