package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetSprintByID(ctx context.Context, id int64) (*models.Sprint, error) {
	builder := sq.Select("id", "title", "description", "type", "target_days", "coins_reward", "is_active", "habit_id", "min_series", "percent_increase", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("sprints").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	sprint := &models.Sprint{}
	var description sql.NullString
	var habitID sql.NullInt64
	var minSeries sql.NullInt64
	var percentIncrease sql.NullInt64
	var sprintType string

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&sprint.ID,
		&sprint.Title,
		&description,
		&sprintType,
		&sprint.TargetDays,
		&sprint.CoinsReward,
		&sprint.IsActive,
		&habitID,
		&minSeries,
		&percentIncrease,
		&sprint.CreatedAt,
		&sprint.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sprint not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get sprint by id: %w", err)
	}

	sprint.Type = models.SprintType(sprintType)
	if description.Valid {
		sprint.Description = description.String
	}
	if habitID.Valid {
		sprint.HabitID = &habitID.Int64
	}
	if minSeries.Valid {
		ms := int(minSeries.Int64)
		sprint.MinSeries = &ms
	}
	if percentIncrease.Valid {
		pi := int(percentIncrease.Int64)
		sprint.PercentIncrease = &pi
	}

	return sprint, nil
}

func (r *Repository) GetAllSprints(ctx context.Context, isActive *bool) ([]*models.Sprint, error) {
	builder := sq.Select("id", "title", "description", "type", "target_days", "coins_reward", "is_active", "habit_id", "min_series", "percent_increase", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("sprints")

	if isActive != nil {
		builder = builder.Where(sq.Eq{"is_active": *isActive})
	}

	builder = builder.OrderBy("created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprints: %w", err)
	}
	defer rows.Close()

	sprints := make([]*models.Sprint, 0)
	for rows.Next() {
		sprint := &models.Sprint{}
		var description sql.NullString
		var habitID sql.NullInt64
		var minSeries sql.NullInt64
		var percentIncrease sql.NullInt64
		var sprintType string

		err := rows.Scan(
			&sprint.ID,
			&sprint.Title,
			&description,
			&sprintType,
			&sprint.TargetDays,
			&sprint.CoinsReward,
			&sprint.IsActive,
			&habitID,
			&minSeries,
			&percentIncrease,
			&sprint.CreatedAt,
			&sprint.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sprint: %w", err)
		}

		sprint.Type = models.SprintType(sprintType)
		if description.Valid {
			sprint.Description = description.String
		}
		if habitID.Valid {
			sprint.HabitID = &habitID.Int64
		}
		if minSeries.Valid {
			ms := int(minSeries.Int64)
			sprint.MinSeries = &ms
		}
		if percentIncrease.Valid {
			pi := int(percentIncrease.Int64)
			sprint.PercentIncrease = &pi
		}

		sprints = append(sprints, sprint)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sprints: %w", err)
	}

	return sprints, nil
}

func (r *Repository) CreateSprint(ctx context.Context, sprint *models.Sprint) (*models.Sprint, error) {
	builder := sq.Insert("sprints").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "description", "type", "target_days", "coins_reward", "is_active", "habit_id", "min_series", "percent_increase").
		Values(sprint.Title, sprint.Description, string(sprint.Type), sprint.TargetDays, sprint.CoinsReward, sprint.IsActive, sprint.HabitID, sprint.MinSeries, sprint.PercentIncrease).
		Suffix("RETURNING id, created_at, updated_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var id int64
	var createdAt, updatedAt sql.NullTime

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create sprint: %w", err)
	}

	sprint.ID = id
	if createdAt.Valid {
		sprint.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		sprint.UpdatedAt = updatedAt.Time
	}

	return sprint, nil
}

func (r *Repository) UpdateSprint(ctx context.Context, sprint *models.Sprint) error {
	builder := sq.Update("sprints").
		PlaceholderFormat(sq.Dollar).
		Set("title", sprint.Title).
		Set("description", sprint.Description).
		Set("type", string(sprint.Type)).
		Set("target_days", sprint.TargetDays).
		Set("coins_reward", sprint.CoinsReward).
		Set("is_active", sprint.IsActive).
		Set("habit_id", sprint.HabitID).
		Set("min_series", sprint.MinSeries).
		Set("percent_increase", sprint.PercentIncrease).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": sprint.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update sprint: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sprint not found")
	}

	return nil
}

func (r *Repository) DeleteSprint(ctx context.Context, id int64) error {
	builder := sq.Delete("sprints").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete sprint: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sprint not found")
	}

	return nil
}

// UserSprintProgress methods

func (r *Repository) GetUserSprintProgress(ctx context.Context, userID int64, sprintID int64) (*models.UserSprintProgress, error) {
	builder := sq.Select("id", "user_id", "sprint_id", "current_days", "is_completed", "completed_at", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("user_sprint_progress").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Eq{"sprint_id": sprintID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	progress := &models.UserSprintProgress{}
	var completedAt sql.NullTime

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.SprintID,
		&progress.CurrentDays,
		&progress.IsCompleted,
		&completedAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user sprint progress: %w", err)
	}

	if completedAt.Valid {
		progress.CompletedAt = &completedAt.Time
	}

	return progress, nil
}

func (r *Repository) GetUserSprintProgresses(ctx context.Context, userID int64) ([]*models.UserSprintProgress, error) {
	builder := sq.Select("id", "user_id", "sprint_id", "current_days", "is_completed", "completed_at", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("user_sprint_progress").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("sprint_id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sprint progresses: %w", err)
	}
	defer rows.Close()

	progresses := make([]*models.UserSprintProgress, 0)
	for rows.Next() {
		progress := &models.UserSprintProgress{}
		var completedAt sql.NullTime

		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.SprintID,
			&progress.CurrentDays,
			&progress.IsCompleted,
			&completedAt,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}

		if completedAt.Valid {
			progress.CompletedAt = &completedAt.Time
		}

		progresses = append(progresses, progress)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating progresses: %w", err)
	}

	return progresses, nil
}

func (r *Repository) CreateOrUpdateUserSprintProgress(ctx context.Context, progress *models.UserSprintProgress) error {
	builder := sq.Insert("user_sprint_progress").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "sprint_id", "current_days", "is_completed", "completed_at").
		Values(progress.UserID, progress.SprintID, progress.CurrentDays, progress.IsCompleted, progress.CompletedAt).
		Suffix("ON CONFLICT (user_id, sprint_id) DO UPDATE SET current_days = EXCLUDED.current_days, is_completed = EXCLUDED.is_completed, completed_at = EXCLUDED.completed_at, updated_at = NOW() RETURNING id, created_at, updated_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build upsert query: %w", err)
	}

	var id int64
	var createdAt, updatedAt sql.NullTime

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return fmt.Errorf("failed to create or update progress: %w", err)
	}

	progress.ID = id
	if createdAt.Valid {
		progress.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		progress.UpdatedAt = updatedAt.Time
	}

	return nil
}

func (r *Repository) ResetAllUserSprintProgresses(ctx context.Context) error {
	builder := sq.Update("user_sprint_progress").
		PlaceholderFormat(sq.Dollar).
		Set("current_days", 0).
		Set("is_completed", false).
		Set("completed_at", nil).
		Set("updated_at", time.Now())

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to reset progresses: %w", err)
	}

	return nil
}
