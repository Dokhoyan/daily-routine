-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Индексы для таблицы habits
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'habits') THEN
        CREATE INDEX IF NOT EXISTS idx_habits_user_id_created_at ON habits(user_id, created_at DESC);
    END IF;

    -- Индексы для таблицы refresh_tokens
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'refresh_tokens') THEN
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at) WHERE revoked_at IS NULL;
    END IF;

    -- Индексы для таблицы token_blacklist
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'token_blacklist') THEN
        CREATE INDEX IF NOT EXISTS idx_token_blacklist_token_hash ON token_blacklist(token_hash);
        CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at);
    END IF;
END $$;
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
