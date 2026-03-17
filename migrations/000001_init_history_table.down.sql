-- Миграция UP: Создание таблицы истории
CREATE TABLE IF NOT EXISTS history (
    id UUID PRIMARY KEY,
    action_name VARCHAR(100) NOT NULL,
    payload TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для ускорения поиска по имени события 
CREATE INDEX IF NOT EXISTS idx_history_action_name ON history(action_name);
