-- Миграция DOWN: Откат изменений
DROP INDEX IF EXISTS idx_history_action_name;
DROP TABLE IF EXISTS history;
