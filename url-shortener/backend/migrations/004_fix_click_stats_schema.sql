-- 1. Меняем тип ip_address с INET на VARCHAR(64) для хранения SHA256-хэша
ALTER TABLE click_stats 
ALTER COLUMN ip_address TYPE VARCHAR(64) USING ip_address::TEXT;

-- 2. (Опционально) Если есть данные в старом формате, можно сконвертировать:
--    Но для новых записей это не нужно, т.к. мы сразу пишем хэш

-- 3. Добавляем индекс для ускорения аналитики (опционально)
CREATE INDEX IF NOT EXISTS idx_click_stats_link_id ON click_stats(link_id);
CREATE INDEX IF NOT EXISTS idx_click_stats_clicked_at ON click_stats(clicked_at);