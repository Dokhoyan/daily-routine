-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_habits_user_id_created_at ON habits(user_id, created_at DESC);

CREATE INDEX idx_refresh_tokens_user_active ON refresh_tokens(user_id, expires_at, created_at DESC) 
WHERE revoked_at IS NULL;

CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

CREATE INDEX idx_token_blacklist_expires_at ON token_blacklist(expires_at);

CREATE INDEX idx_habits_user_type_active ON habits(user_id, type, is_active, created_at DESC);

CREATE INDEX idx_habits_is_active ON habits(is_active) WHERE is_active = TRUE;

CREATE INDEX idx_refresh_tokens_token_active ON refresh_tokens(token) 
WHERE revoked_at IS NULL;

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- Индекс для token_blacklist: оптимизация проверки токена в черном списке
-- Используется в IsTokenBlacklisted (WHERE token_hash = ? AND expires_at > NOW())
-- Примечание: условие expires_at > NOW() проверяется в запросе, не в индексе
CREATE INDEX idx_token_blacklist_hash_active ON token_blacklist(token_hash, expires_at);

CREATE INDEX idx_token_blacklist_user_id ON token_blacklist(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_habits_user_id_created_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_active;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_token_blacklist_expires_at;
DROP INDEX IF EXISTS idx_habits_user_type_active;
DROP INDEX IF EXISTS idx_habits_is_active;
DROP INDEX IF EXISTS idx_refresh_tokens_token_active;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_token_blacklist_hash_active;
DROP INDEX IF EXISTS idx_token_blacklist_user_id;

-- +goose StatementEnd
