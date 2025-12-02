package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetHabitByID(ctx context.Context, id int64) (*models.Habit, error) {
	builder := sq.Select("id", "user_id", "title", "type", "unit", "value", "is_active", "is_done", "is_beneficial", "series", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("habits").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	habit := &models.Habit{}
	var habitType string

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Title,
		&habitType,
		&habit.Unit,
		&habit.Value,
		&habit.IsActive,
		&habit.IsDone,
		&habit.IsBeneficial,
		&habit.Series,
		&habit.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("habit not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get habit by id: %w", err)
	}

	habit.Type = models.HabitType(habitType)
	return habit, nil
}

func (r *Repository) GetHabitsByUserID(ctx context.Context, userID int64) ([]*models.Habit, error) {
	builder := sq.Select("id", "user_id", "title", "type", "unit", "value", "is_active", "is_done", "is_beneficial", "series", "created_at").
		PlaceholderFormat(sq.Dollar).
		From("habits").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits by user id: %w", err)
	}
	defer rows.Close()

	habits := make([]*models.Habit, 0)
	for rows.Next() {
		habit := &models.Habit{}
		var habitType string

		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Title,
			&habitType,
			&habit.Unit,
			&habit.Value,
			&habit.IsActive,
			&habit.IsDone,
			&habit.IsBeneficial,
			&habit.Series,
			&habit.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}

		habit.Type = models.HabitType(habitType)
		habits = append(habits, habit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

func (r *Repository) CreateHabit(ctx context.Context, habit *models.Habit) (*models.Habit, error) {
	builder := sq.Insert("habits").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "title", "type", "unit", "value", "is_active", "is_done", "is_beneficial", "series").
		Values(habit.UserID, habit.Title, string(habit.Type), habit.Unit, habit.Value, habit.IsActive, habit.IsDone, habit.IsBeneficial, habit.Series).
		Suffix("RETURNING id, created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var id int64
	var createdAt sql.NullTime

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create habit: %w", err)
	}

	habit.ID = id
	if createdAt.Valid {
		habit.CreatedAt = createdAt.Time
	}

	return habit, nil
}

func (r *Repository) UpdateHabit(ctx context.Context, habit *models.Habit) error {
	builder := sq.Update("habits").
		PlaceholderFormat(sq.Dollar).
		Set("title", habit.Title).
		Set("type", string(habit.Type)).
		Set("unit", habit.Unit).
		Set("value", habit.Value).
		Set("is_active", habit.IsActive).
		Set("is_done", habit.IsDone).
		Set("is_beneficial", habit.IsBeneficial).
		Set("series", habit.Series).
		Where(sq.Eq{"id": habit.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("habit not found")
	}

	return nil
}

func (r *Repository) DeleteHabit(ctx context.Context, id int64) error {
	builder := sq.Delete("habits").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("habit not found")
	}

	return nil
}

func (r *Repository) DeleteHabitsByUserID(ctx context.Context, userID int64) error {
	builder := sq.Delete("habits").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_id": userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete habits by user id: %w", err)
	}

	return nil
}
