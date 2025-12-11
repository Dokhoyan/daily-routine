-- +goose Up
-- +goose StatementBegin
ALTER TABLE habits
    ADD COLUMN IF NOT EXISTS current_value INTEGER DEFAULT 0 CHECK (current_value <= value);

UPDATE habits
SET current_value = LEAST(COALESCE(current_value, 0), value)
WHERE current_value IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE habits
    DROP COLUMN IF EXISTS current_value;
-- +goose StatementEnd
