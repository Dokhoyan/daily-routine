package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func (r *Repository) GetSettingsByUserID(ctx context.Context, userID int64) (*models.UserSettings, error) {
	builder := sq.Select("user_id", "timezone", "do_not_disturb", "notify_times").
		PlaceholderFormat(sq.Dollar).
		From("user_settings").
		Where(sq.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	settings := &models.UserSettings{}
	var notifyTimes pq.StringArray

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&settings.UserID,
		&settings.Timezone,
		&settings.DoNotDisturb,
		&notifyTimes,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.UserSettings{
				UserID:       userID,
				Timezone:     "UTC",
				DoNotDisturb: false,
				NotifyTimes:  []string{},
			}, nil
		}
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}

	settings.NotifyTimes = []string(notifyTimes)
	return settings, nil
}

func (r *Repository) CreateSettings(ctx context.Context, settings *models.UserSettings) error {
	builder := sq.Insert("user_settings").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "timezone", "do_not_disturb", "notify_times").
		Values(settings.UserID, settings.Timezone, settings.DoNotDisturb, pq.Array(settings.NotifyTimes))

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create user settings: %w", err)
	}

	return nil
}

func (r *Repository) UpdateSettings(ctx context.Context, settings *models.UserSettings) error {
	builder := sq.Update("user_settings").
		PlaceholderFormat(sq.Dollar).
		Set("timezone", settings.Timezone).
		Set("do_not_disturb", settings.DoNotDisturb).
		Set("notify_times", pq.Array(settings.NotifyTimes)).
		Where(sq.Eq{"user_id": settings.UserID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user settings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return r.CreateSettings(ctx, settings)
	}

	return nil
}
