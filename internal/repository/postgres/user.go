package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dokhoyan/daily-routine/internal/models"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	builder := sq.Select("id", "username", "photo_url", "auth_date", "tokentg").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	user := &models.User{}
	var authDate sql.NullTime

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.PhotoURL,
		&authDate,
		&user.TokenTG,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if authDate.Valid {
		user.AuthDate = authDate.Time
	}

	return user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "username", "photo_url", "auth_date", "tokentg").
		Values(user.ID, user.Username, user.PhotoURL, user.AuthDate, user.TokenTG)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	builder := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("username", user.Username).
		Set("photo_url", user.PhotoURL).
		Set("auth_date", user.AuthDate).
		Set("tokentg", user.TokenTG).
		Where(sq.Eq{"id": user.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	builder := sq.Select("id", "username", "photo_url", "auth_date", "tokentg").
		PlaceholderFormat(sq.Dollar).
		From("users").
		OrderBy("id")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		var authDate sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PhotoURL,
			&authDate,
			&user.TokenTG,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if authDate.Valid {
			user.AuthDate = authDate.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
