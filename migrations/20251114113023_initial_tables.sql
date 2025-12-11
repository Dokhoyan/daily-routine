-- +goose Up
-- +goose StatementBegin
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

CREATE TYPE habit_format AS ENUM ('time', 'count', 'binary');
CREATE TYPE habit_type AS ENUM ('beneficial', 'harmful');

CREATE TABLE habits (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES users(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    format          habit_format NOT NULL,
    unit            TEXT,
    value           INTEGER NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE,
    is_done         BOOLEAN DEFAULT FALSE,
    type            habit_type NOT NULL DEFAULT 'beneficial',
    series          INTEGER DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS token_issuance_log;
DROP TABLE IF EXISTS token_blacklist;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS habits;
DROP TYPE IF EXISTS habit_type;
DROP TYPE IF EXISTS habit_format;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

