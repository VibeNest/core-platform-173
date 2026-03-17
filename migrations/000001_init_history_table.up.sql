
-- Миграция UP: Создание таблицы истории
CREATE TABLE IF NOT EXISTS history (
    id UUID PRIMARY KEY,
    source VARCHAR(10) NOT NULL,
    destination VARCHAR(10) NOT NULL,
    original TEXT NOT NULL,
    translation TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- Индекс для ускорения поиска по имени события 