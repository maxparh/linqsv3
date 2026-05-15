package repository

import (
	"context"
	"database/sql"
	"url-shortener/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.AnalyticsSession) error
	GetBySessionID(ctx context.Context, sessionID string) (*domain.AnalyticsSession, error)
	Update(ctx context.Context, session *domain.AnalyticsSession) error
	GetAnalyticsByLinkID(ctx context.Context, linkID int, days int) ([]*domain.AnalyticsSession, error)
	RecordSession(ctx context.Context, linkID int, sessionID string, ipHash string, countryCode string, deviceType string, browserName string) error

	GetOverviewStats(ctx context.Context, userID int, days int) (*domain.AnalyticsOverview, error)
	GetClicksOverTime(ctx context.Context, userID int, days int) (*domain.ClicksOverTime, error)
	GetTopLocations(ctx context.Context, userID int, days int, limit int) ([]*domain.LocationStat, error)
	GetDeviceStats(ctx context.Context, userID int, days int) ([]*domain.DeviceStat, error)
}

type postgresSessionRepo struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &postgresSessionRepo{db: db}
}

func (r *postgresSessionRepo) Create(ctx context.Context, session *domain.AnalyticsSession) error {
	query := `INSERT INTO analytics_sessions 
        (session_id, link_id, user_ip_hash, country_code, device_type, browser_name, 
         started_at, last_activity_at, page_views, total_time_on_site, is_bounce)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		session.SessionID,
		session.LinkID,
		session.UserIPHash,
		session.CountryCode,
		session.DeviceType,
		session.BrowserName,
		session.StartedAt,
		session.LastActivityAt,
		session.PageViews,
		session.TotalTimeOnSite,
		session.IsBounce,
	).Scan(&session.ID)
}

func (r *postgresSessionRepo) GetBySessionID(ctx context.Context, sessionID string) (*domain.AnalyticsSession, error) {
	var session domain.AnalyticsSession
	query := `SELECT id, session_id, link_id, user_ip_hash, started_at, last_activity_at,
                     page_views, total_time_on_site, is_bounce, country_code, device_type, browser_name
              FROM analytics_sessions WHERE session_id = $1`

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &session.SessionID, &session.LinkID, &session.UserIPHash,
		&session.StartedAt, &session.LastActivityAt, &session.PageViews,
		&session.TotalTimeOnSite, &session.IsBounce, &session.CountryCode,
		&session.DeviceType, &session.BrowserName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *postgresSessionRepo) Update(ctx context.Context, session *domain.AnalyticsSession) error {
	query := `UPDATE analytics_sessions 
        SET page_views = $1, total_time_on_site = $2, is_bounce = $3, 
            last_activity_at = $4
        WHERE session_id = $5`

	_, err := r.db.ExecContext(ctx, query,
		session.PageViews,
		session.TotalTimeOnSite,
		session.IsBounce,
		session.LastActivityAt,
		session.SessionID,
	)

	return err
}

func (r *postgresSessionRepo) GetAnalyticsByLinkID(ctx context.Context, linkID int, days int) ([]*domain.AnalyticsSession, error) {
	query := `SELECT id, session_id, link_id, user_ip_hash, started_at, last_activity_at,
                     page_views, total_time_on_site, is_bounce, country_code, device_type, browser_name
              FROM analytics_sessions 
              WHERE link_id = $1 AND started_at >= NOW() - INTERVAL '1 day' * $2
              ORDER BY started_at DESC`

	rows, err := r.db.QueryContext(ctx, query, linkID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*domain.AnalyticsSession
	for rows.Next() {
		var session domain.AnalyticsSession
		err := rows.Scan(
			&session.ID, &session.SessionID, &session.LinkID, &session.UserIPHash,
			&session.StartedAt, &session.LastActivityAt, &session.PageViews,
			&session.TotalTimeOnSite, &session.IsBounce, &session.CountryCode,
			&session.DeviceType, &session.BrowserName,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, rows.Err()
}

func (r *postgresSessionRepo) RecordSession(ctx context.Context, linkID int, sessionID string, ipHash string, countryCode string, deviceType string, browserName string) error {
	query := `INSERT INTO analytics_sessions 
        (link_id, session_id, user_ip_hash, country_code, device_type, browser_name, is_bounce, page_views, started_at, last_activity_at) 
        VALUES ($1, $2, $3, $4, $5, $6, true, 1, NOW(), NOW())
        ON CONFLICT (session_id) DO UPDATE SET
            link_id = EXCLUDED.link_id,
            page_views = analytics_sessions.page_views + 1,
            last_activity_at = NOW(),
            is_bounce = false`

	_, err := r.db.ExecContext(ctx, query, linkID, sessionID, ipHash, countryCode, deviceType, browserName)
	return err
}

func (r *postgresSessionRepo) GetOverviewStats(ctx context.Context, userID int, days int) (*domain.AnalyticsOverview, error) {
	var o domain.AnalyticsOverview

	query := `
        SELECT 
            COALESCE(SUM(s.page_views), 0)::bigint as total_clicks,
            COUNT(DISTINCT s.session_id)::bigint as unique_clicks,
            COALESCE(AVG(s.total_time_on_site)::float, 0) as avg_time_on_site,
            CASE 
                WHEN COUNT(DISTINCT s.session_id) = 0 THEN 0 
                ELSE COUNT(DISTINCT CASE WHEN s.is_bounce THEN s.session_id END) * 100.0 / COUNT(DISTINCT s.session_id) 
            END as bounce_rate
        FROM analytics_sessions s
        JOIN links l ON s.link_id = l.id
        WHERE l.user_id = $1 AND s.started_at >= NOW() - make_interval(days => $2)
    `
	err := r.db.QueryRowContext(ctx, query, userID, days).Scan(&o.TotalClicks, &o.UniqueClicks, &o.AvgTimeOnSite, &o.BounceRate)
	return &o, err
}

func (r *postgresSessionRepo) GetClicksOverTime(ctx context.Context, userID int, days int) (*domain.ClicksOverTime, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT to_char(s.started_at, 'DD.MM') as day, SUM(s.page_views) as clicks
        FROM analytics_sessions s
        JOIN links l ON s.link_id = l.id
        WHERE l.user_id = $1 AND s.started_at >= NOW() - make_interval(days => $2)
        GROUP BY day ORDER BY day
    `, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res domain.ClicksOverTime
	for rows.Next() {
		var day string
		var clicks int
		if err := rows.Scan(&day, &clicks); err != nil {
			return nil, err
		}
		res.Labels = append(res.Labels, day)
		res.Values = append(res.Values, clicks)
	}
	return &res, nil
}

func (r *postgresSessionRepo) GetTopLocations(ctx context.Context, userID int, days int, limit int) ([]*domain.LocationStat, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT s.country_code, COUNT(*), ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2)
        FROM analytics_sessions s JOIN links l ON s.link_id = l.id
        WHERE l.user_id = $1 AND s.started_at >= NOW() - make_interval(days => $2) AND s.country_code != ''
        GROUP BY s.country_code ORDER BY count DESC LIMIT $3
    `, userID, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := map[string]string{"RU": "Россия", "BY": "Беларусь", "DE": "Германия", "FI": "Финляндия", "FR": "Франция", "US": "США", "UA": "Украина", "KZ": "Казахстан"}
	var stats []*domain.LocationStat
	for rows.Next() {
		var code string
		var cnt int
		var pct float64
		rows.Scan(&code, &cnt, &pct)
		stats = append(stats, &domain.LocationStat{Country: names[code], CountryCode: code, Percent: pct})
	}
	return stats, nil
}

func (r *postgresSessionRepo) GetDeviceStats(ctx context.Context, userID int, days int) ([]*domain.DeviceStat, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT s.device_type, COUNT(*), ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2)
        FROM analytics_sessions s JOIN links l ON s.link_id = l.id
        WHERE l.user_id = $1 AND s.started_at >= NOW() - make_interval(days => $2) AND s.device_type != ''
        GROUP BY s.device_type ORDER BY count DESC
    `, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := map[string]string{"desktop": "Компьютеры", "mobile": "Смартфоны", "tablet": "Планшеты"}
	colors := map[string]string{"desktop": "#014751", "mobile": "#14b8a6", "tablet": "#10b981"}
	var stats []*domain.DeviceStat
	for rows.Next() {
		var dev string
		var cnt int
		var pct float64
		rows.Scan(&dev, &cnt, &pct)
		stats = append(stats, &domain.DeviceStat{Name: names[dev], Value: cnt, Percent: pct, Color: colors[dev]})
	}
	return stats, nil
}
