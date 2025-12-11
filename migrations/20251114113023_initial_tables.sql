-- +goose Up
-- +goose StatementBegin
-- Users and settings
CREATE TABLE users (
    id          BIGINT PRIMARY KEY,
    username    TEXT,
    first_name  TEXT,
    photo_url   TEXT,
    auth_date   TIMESTAMP,
    tokentg     TEXT
);

CREATE TABLE user_settings (
    user_id         BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    timezone        TEXT DEFAULT 'UTC',
    do_not_disturb  BOOLEAN DEFAULT FALSE,
    notify_times    TEXT[] DEFAULT '{}'
);

-- Habit enums
CREATE TYPE habit_format AS ENUM ('time', 'count', 'binary');
CREATE TYPE habit_type AS ENUM ('beneficial', 'harmful');

-- Habits
CREATE TABLE habits (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES users(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    format          habit_format NOT NULL,
    unit            TEXT,
    value           INTEGER NOT NULL,
    current_value   INTEGER DEFAULT 0 CHECK (current_value <= value),
    is_active       BOOLEAN DEFAULT TRUE,
    is_done         BOOLEAN DEFAULT FALSE,
    type            habit_type NOT NULL DEFAULT 'beneficial',
    series          INTEGER DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Tokens
CREATE TABLE refresh_tokens (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           TEXT NOT NULL UNIQUE,
    expires_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    revoked_at      TIMESTAMP,
    device_info     TEXT,
    ip_address      TEXT
);

CREATE TABLE token_blacklist (
    id              BIGSERIAL PRIMARY KEY,
    token_hash      TEXT NOT NULL UNIQUE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at      TIMESTAMP NOT NULL,
    revoked_at      TIMESTAMP DEFAULT NOW(),
    reason          TEXT
);

CREATE TABLE token_issuance_log (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_type      TEXT NOT NULL,
    action          TEXT NOT NULL,
    device_info     TEXT,
    ip_address      TEXT,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_habits_user_id_created_at ON habits(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_token_blacklist_token_hash ON token_blacklist(token_hash);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at);
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
DROP TABLE IF EXISTS token_issuance_log;
DROP TABLE IF EXISTS token_blacklist;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS habits;
DROP TYPE IF EXISTS habit_type;
DROP TYPE IF EXISTS habit_format;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd