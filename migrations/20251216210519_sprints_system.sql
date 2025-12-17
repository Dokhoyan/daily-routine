-- +goose Up
-- +goose StatementBegin

-- Добавляем поле is_admin в таблицу users
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT FALSE;

-- Создаем тип для видов спринтов
CREATE TYPE sprint_type AS ENUM ('habit_series', 'all_habits', 'habit_increase');

-- Обновляем таблицу sprints с полной структурой
ALTER TABLE sprints DROP COLUMN IF EXISTS title;
ALTER TABLE sprints DROP COLUMN IF EXISTS value;
ALTER TABLE sprints DROP COLUMN IF EXISTS current_value;
ALTER TABLE sprints DROP COLUMN IF EXISTS coins;
ALTER TABLE sprints DROP COLUMN IF EXISTS created_at;

ALTER TABLE sprints ADD COLUMN IF NOT EXISTS title TEXT NOT NULL;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS type sprint_type NOT NULL;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS target_days INTEGER NOT NULL;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS coins_reward INTEGER NOT NULL DEFAULT 0;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Параметры для разных типов спринтов
-- Для habit_series: habit_id - ID привычки, min_series - минимальная серия
-- Для all_habits: нет параметров
-- Для habit_increase: habit_id - ID привычки, percent_increase - процент увеличения
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS habit_id BIGINT REFERENCES habits(id) ON DELETE CASCADE;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS min_series INTEGER;
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS percent_increase INTEGER;

-- Таблица для отслеживания прогресса пользователей по спринтам
CREATE TABLE IF NOT EXISTS user_sprint_progress (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sprint_id       BIGINT NOT NULL REFERENCES sprints(id) ON DELETE CASCADE,
    current_days    INTEGER DEFAULT 0,
    is_completed    BOOLEAN DEFAULT FALSE,
    completed_at     TIMESTAMP,
    week_start      DATE NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, sprint_id, week_start)
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_user_sprint_progress_user_sprint ON user_sprint_progress(user_id, sprint_id);
CREATE INDEX IF NOT EXISTS idx_user_sprint_progress_week_start ON user_sprint_progress(week_start);
CREATE INDEX IF NOT EXISTS idx_sprints_type_active ON sprints(type, is_active);
CREATE INDEX IF NOT EXISTS idx_sprints_created_by ON sprints(created_by);

-- Создаем админского пользователя (ID: 1, username: admin, password будет через Telegram)
-- Пароль для Telegram бота не нужен, но создадим пользователя с is_admin = true
INSERT INTO users (id, username, first_name, is_admin, coins) 
VALUES (1, 'admin', 'Administrator', TRUE, 0)
ON CONFLICT (id) DO UPDATE SET is_admin = TRUE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_sprints_created_by;
DROP INDEX IF EXISTS idx_sprints_type_active;
DROP INDEX IF EXISTS idx_user_sprint_progress_week_start;
DROP INDEX IF EXISTS idx_user_sprint_progress_user_sprint;

DROP TABLE IF EXISTS user_sprint_progress;

ALTER TABLE sprints DROP COLUMN IF EXISTS percent_increase;
ALTER TABLE sprints DROP COLUMN IF EXISTS min_series;
ALTER TABLE sprints DROP COLUMN IF EXISTS habit_id;
ALTER TABLE sprints DROP COLUMN IF EXISTS updated_at;
ALTER TABLE sprints DROP COLUMN IF EXISTS created_at;
ALTER TABLE sprints DROP COLUMN IF EXISTS created_by;
ALTER TABLE sprints DROP COLUMN IF EXISTS is_active;
ALTER TABLE sprints DROP COLUMN IF EXISTS coins_reward;
ALTER TABLE sprints DROP COLUMN IF EXISTS target_days;
ALTER TABLE sprints DROP COLUMN IF EXISTS type;
ALTER TABLE sprints DROP COLUMN IF EXISTS description;
ALTER TABLE sprints DROP COLUMN IF EXISTS title;

ALTER TABLE sprints ADD COLUMN created_at TIMESTAMP DEFAULT NOW();
ALTER TABLE sprints ADD COLUMN coins INTEGER NOT NULL;
ALTER TABLE sprints ADD COLUMN current_value INTEGER DEFAULT 0 CHECK (current_value <= value);
ALTER TABLE sprints ADD COLUMN value INTEGER NOT NULL;
ALTER TABLE sprints ADD COLUMN title TEXT NOT NULL;

DROP TYPE IF EXISTS sprint_type;

ALTER TABLE users DROP COLUMN IF EXISTS is_admin;

-- +goose StatementEnd
