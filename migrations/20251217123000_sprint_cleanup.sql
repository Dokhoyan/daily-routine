-- +goose Up
-- +goose StatementBegin

-- Удаляем поле created_by из спринтов (больше не храним автора)
DROP INDEX IF EXISTS idx_sprints_created_by;
ALTER TABLE sprints DROP COLUMN IF EXISTS created_by;

-- Прогресс храним без week_start, т.к. еженедельный сброс идет по крону
ALTER TABLE user_sprint_progress DROP CONSTRAINT IF EXISTS user_sprint_progress_user_id_sprint_id_week_start_key;
DROP INDEX IF EXISTS idx_user_sprint_progress_week_start;
ALTER TABLE user_sprint_progress DROP COLUMN IF EXISTS week_start;

-- Новый уникальный ключ по пользователю и спринту
ALTER TABLE user_sprint_progress
  ADD CONSTRAINT user_sprint_progress_unique UNIQUE (user_id, sprint_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Возвращаем поле week_start и индекс/уникальность
ALTER TABLE user_sprint_progress DROP CONSTRAINT IF EXISTS user_sprint_progress_unique;
ALTER TABLE user_sprint_progress ADD COLUMN IF NOT EXISTS week_start DATE NOT NULL DEFAULT CURRENT_DATE;
CREATE UNIQUE INDEX IF NOT EXISTS user_sprint_progress_user_id_sprint_id_week_start_key
  ON user_sprint_progress(user_id, sprint_id, week_start);
CREATE INDEX IF NOT EXISTS idx_user_sprint_progress_week_start ON user_sprint_progress(week_start);

-- Возвращаем поле created_by и индекс
ALTER TABLE sprints ADD COLUMN IF NOT EXISTS created_by BIGINT REFERENCES users(id);
CREATE INDEX IF NOT EXISTS idx_sprints_created_by ON sprints(created_by);

-- +goose StatementEnd


