-- +goose Up
-- +goose StatementBegin
-- Индексы для таблицы habits
-- Используется в GetHabitsByUserID: WHERE user_id = ? ORDER BY created_at DESC
-- Составной индекс покрывает и WHERE user_id = ? (левая часть индекса)
CREATE INDEX idx_habits_user_id_created_at ON habits(user_id, created_at DESC);

-- Индексы для таблицы refresh_tokens
-- Используется в GetRefreshTokenByHash и RevokeRefreshToken: WHERE token = ?
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
-- Используется в RevokeAllUserTokens, GetActiveTokensCount, GetActiveTokens: WHERE user_id = ?
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
-- Используется в DeleteExpiredTokens: WHERE expires_at < NOW()
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
-- Используется в GetActiveTokensCount, GetActiveTokens: WHERE revoked_at IS NULL
-- Частичный индекс для оптимизации запросов по активным токенам
CREATE INDEX idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at) WHERE revoked_at IS NULL;

-- Индексы для таблицы token_blacklist
-- Используется в IsTokenBlacklisted: WHERE token_hash = ? AND expires_at > NOW()
CREATE INDEX idx_token_blacklist_token_hash ON token_blacklist(token_hash);
-- Используется в DeleteExpiredBlacklistEntries: WHERE expires_at < NOW()
CREATE INDEX idx_token_blacklist_expires_at ON token_blacklist(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_token_blacklist_expires_at;
DROP INDEX IF EXISTS idx_token_blacklist_token_hash;
DROP INDEX IF EXISTS idx_refresh_tokens_revoked_at;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_habits_user_id_created_at;
-- +goose StatementEnd
