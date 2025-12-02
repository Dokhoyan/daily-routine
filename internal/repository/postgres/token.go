package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/models"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	builder := sq.Insert("refresh_tokens").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "token", "expires_at", "device_info", "ip_address").
		Values(token.UserID, token.Token, token.ExpiresAt, token.DeviceInfo, token.IPAddress)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func (r *Repository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	builder := sq.Select("id", "user_id", "token", "expires_at", "created_at", "revoked_at", "device_info", "ip_address").
		PlaceholderFormat(sq.Dollar).
		From("refresh_tokens").
		Where(sq.Eq{"token": tokenHash}).
		Where(sq.Expr("revoked_at IS NULL"))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	token := &models.RefreshToken{}
	var revokedAt sql.NullTime
	var deviceInfo, ipAddress sql.NullString

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&revokedAt,
		&deviceInfo,
		&ipAddress,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found or revoked: %w", err)
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	if deviceInfo.Valid {
		token.DeviceInfo = &deviceInfo.String
	}
	if ipAddress.Valid {
		token.IPAddress = &ipAddress.String
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	return token, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	builder := sq.Update("refresh_tokens").
		PlaceholderFormat(sq.Dollar).
		Set("revoked_at", time.Now()).
		Where(sq.Eq{"token": tokenHash}).
		Where(sq.Expr("revoked_at IS NULL"))

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found or already revoked")
	}

	return nil
}

func (r *Repository) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	builder := sq.Update("refresh_tokens").
		PlaceholderFormat(sq.Dollar).
		Set("revoked_at", time.Now()).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Expr("revoked_at IS NULL"))

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}

	return nil
}

func (r *Repository) GetActiveTokensCount(ctx context.Context, userID int64) (int, error) {
	builder := sq.Select("COUNT(*)").
		PlaceholderFormat(sq.Dollar).
		From("refresh_tokens").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Expr("revoked_at IS NULL")).
		Where(sq.Expr("expires_at > NOW()"))

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var count int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get active tokens count: %w", err)
	}

	return count, nil
}

func (r *Repository) GetActiveTokens(ctx context.Context, userID int64) ([]*models.RefreshToken, error) {
	builder := sq.Select("id", "user_id", "token", "expires_at", "created_at", "revoked_at", "device_info", "ip_address").
		PlaceholderFormat(sq.Dollar).
		From("refresh_tokens").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Expr("revoked_at IS NULL")).
		Where(sq.Expr("expires_at > NOW()")).
		OrderBy("created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get active tokens: %w", err)
	}
	defer rows.Close()

	tokens := make([]*models.RefreshToken, 0)
	for rows.Next() {
		token := &models.RefreshToken{}
		var revokedAt sql.NullTime
		var deviceInfo, ipAddress sql.NullString

		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.CreatedAt,
			&revokedAt,
			&deviceInfo,
			&ipAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}

		if revokedAt.Valid {
			token.RevokedAt = &revokedAt.Time
		}
		if deviceInfo.Valid {
			token.DeviceInfo = &deviceInfo.String
		}
		if ipAddress.Valid {
			token.IPAddress = &ipAddress.String
		}

		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tokens: %w", err)
	}

	return tokens, nil
}

func (r *Repository) DeleteExpiredTokens(ctx context.Context) error {
	builder := sq.Delete("refresh_tokens").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Expr("expires_at < NOW()"))

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	return nil
}

func (r *Repository) LogTokenAction(ctx context.Context, log *models.TokenLog) error {
	builder := sq.Insert("token_issuance_log").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "token_type", "action", "device_info", "ip_address").
		Values(log.UserID, log.TokenType, log.Action, log.DeviceInfo, log.IPAddress)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to log token action: %w", err)
	}

	return nil
}

func (r *Repository) AddToBlacklist(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time, reason *string) error {
	builder := sq.Insert("token_blacklist").
		PlaceholderFormat(sq.Dollar).
		Columns("token_hash", "user_id", "expires_at", "reason").
		Values(tokenHash, userID, expiresAt, reason)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	return nil
}

func (r *Repository) IsTokenBlacklisted(ctx context.Context, tokenHash string) (bool, error) {
	builder := sq.Select("COUNT(*)").
		PlaceholderFormat(sq.Dollar).
		From("token_blacklist").
		Where(sq.Eq{"token_hash": tokenHash}).
		Where(sq.Expr("expires_at > NOW()"))

	query, args, err := builder.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build select query: %w", err)
	}

	var count int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) DeleteExpiredBlacklistEntries(ctx context.Context) error {
	builder := sq.Delete("token_blacklist").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Expr("expires_at < NOW()"))

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete expired blacklist entries: %w", err)
	}

	return nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
