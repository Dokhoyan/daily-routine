-- +goose Up
-- +goose StatementBegin
-- 1) Форматы привычек (time/count/binary) переносим в отдельный enum/колонку format
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'habit_format' AND n.nspname = 'public'
    ) THEN
        CREATE TYPE habit_format AS ENUM ('time', 'count', 'binary');
    END IF;
END$$;

ALTER TABLE habits
    ADD COLUMN IF NOT EXISTS format habit_format;

-- заполняем format из старого столбца type (enum со значениями time/count/binary)
UPDATE habits
SET format = CASE
    WHEN format IS NOT NULL THEN format
    WHEN type::text = 'time' THEN 'time'::habit_format
    WHEN type::text = 'count' THEN 'count'::habit_format
    WHEN type::text = 'binary' THEN 'binary'::habit_format
    ELSE 'count'::habit_format
END
WHERE format IS NULL;

ALTER TABLE habits ALTER COLUMN format SET NOT NULL;
ALTER TABLE habits ALTER COLUMN format SET DEFAULT 'count';

-- 2) Новый enum для полезных/вредных привычек и колонка type
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'habit_type_new' AND n.nspname = 'public'
    ) THEN
        CREATE TYPE habit_type_new AS ENUM ('beneficial', 'harmful');
    END IF;
END$$;

ALTER TABLE habits
    ADD COLUMN IF NOT EXISTS type_new habit_type_new DEFAULT 'beneficial'::habit_type_new;

UPDATE habits
SET type_new = CASE
    WHEN is_beneficial = FALSE THEN 'harmful'::habit_type_new
    ELSE 'beneficial'::habit_type_new
END;

ALTER TABLE habits ALTER COLUMN type_new SET NOT NULL;
ALTER TABLE habits ALTER COLUMN type_new SET DEFAULT 'beneficial';

-- удаляем старые столбцы и типы, приводим к новой схеме
ALTER TABLE habits DROP COLUMN IF EXISTS is_beneficial;
ALTER TABLE habits DROP COLUMN IF EXISTS type;
ALTER TABLE habits RENAME COLUMN type_new TO type;

-- сносим старый enum habit_type, если он остался, и переименовываем новый
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'habit_type' AND n.nspname = 'public'
    ) THEN
        DROP TYPE habit_type;
    END IF;
END$$;

ALTER TYPE habit_type_new RENAME TO habit_type;

-- 3) Добавляем first_name пользователю (idempotent)
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS first_name TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Откат к старой схеме (best-effort)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'habit_type_old' AND n.nspname = 'public'
    ) THEN
        CREATE TYPE habit_type_old AS ENUM ('time', 'count', 'binary');
    END IF;
END$$;

ALTER TABLE habits
    ADD COLUMN IF NOT EXISTS type_old habit_type_old;

-- восстанавливаем type_old из format
UPDATE habits
SET type_old = CASE
    WHEN format::text = 'time' THEN 'time'::habit_type_old
    WHEN format::text = 'count' THEN 'count'::habit_type_old
    WHEN format::text = 'binary' THEN 'binary'::habit_type_old
    ELSE 'count'::habit_type_old
END
WHERE type_old IS NULL;

ALTER TABLE habits
    ADD COLUMN IF NOT EXISTS is_beneficial BOOLEAN DEFAULT TRUE;

UPDATE habits
SET is_beneficial = CASE
    WHEN type::text = 'harmful' THEN FALSE
    ELSE TRUE
END;

ALTER TABLE habits DROP COLUMN IF EXISTS type;
ALTER TABLE habits RENAME COLUMN type_old TO type;

-- возвращаем старый enum habit_type, если его нет
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'habit_type' AND n.nspname = 'public'
    ) THEN
        DROP TYPE habit_type;
    END IF;
END$$;

ALTER TYPE habit_type_old RENAME TO habit_type;

-- удаляем новый формат и колонку
ALTER TABLE habits DROP COLUMN IF EXISTS format;
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'habit_format' AND n.nspname = 'public'
    ) THEN
        DROP TYPE habit_format;
    END IF;
END$$;

-- оставляем first_name, так как откат может затронуть данные; удаляем опционально
-- ALTER TABLE users DROP COLUMN IF EXISTS first_name;
-- +goose StatementEnd
