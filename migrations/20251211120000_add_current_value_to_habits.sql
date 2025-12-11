-- Add current_value column for partial habit completion
ALTER TABLE habits
    ADD COLUMN IF NOT EXISTS current_value INTEGER DEFAULT 0 CHECK (current_value <= value);

-- Backfill existing rows to a safe value
UPDATE habits
SET current_value = LEAST(COALESCE(current_value, 0), value)
WHERE current_value IS NULL;
