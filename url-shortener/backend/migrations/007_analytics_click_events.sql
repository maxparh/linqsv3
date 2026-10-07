ALTER TABLE analytics_sessions ADD COLUMN event_based BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE analytics_click_events (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL REFERENCES analytics_sessions(session_id) ON DELETE CASCADE,
    link_id INTEGER NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    country_code CHAR(2),
    device_type VARCHAR(50)
);
CREATE INDEX idx_click_events_link_time ON analytics_click_events(link_id, clicked_at);
CREATE INDEX idx_click_events_session ON analytics_click_events(session_id);

-- Legacy totals retain their original dates; individual historical dates are unknown.
CREATE VIEW analytics_click_totals AS
SELECT session_id, link_id, clicked_at, country_code, device_type, 1 AS clicks
FROM analytics_click_events
UNION ALL
SELECT session_id, link_id, started_at, country_code, device_type, page_views
FROM analytics_sessions WHERE NOT event_based;
