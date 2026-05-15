-- Migration 002: Add analytics_sessions table for proper session tracking

-- Таблица сессий для аналитики
CREATE TABLE IF NOT EXISTS analytics_sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL UNIQUE,
    link_id INTEGER REFERENCES links(id) ON DELETE CASCADE,
    user_ip_hash VARCHAR(64),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    page_views INTEGER DEFAULT 1,
    total_time_on_site INTEGER DEFAULT 0,
    is_bounce BOOLEAN DEFAULT TRUE,
    country_code CHAR(2),
    device_type VARCHAR(50),
    browser_name VARCHAR(50)
);

-- Индексы для ускорения запросов
CREATE INDEX IF NOT EXISTS idx_sessions_link_id ON analytics_sessions(link_id);
CREATE INDEX IF NOT EXISTS idx_sessions_session_id ON analytics_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_sessions_started_at ON analytics_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_sessions_last_activity ON analytics_sessions(last_activity_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_ip ON analytics_sessions(user_ip_hash);

-- Комментарий к таблице
COMMENT ON TABLE analytics_sessions IS 'Sessions for analytics tracking - one row per user session per link';
COMMENT ON COLUMN analytics_sessions.is_bounce IS 'TRUE if session had only 1 page view (bounce)';
COMMENT ON COLUMN analytics_sessions.total_time_on_site IS 'Time spent on site in seconds';