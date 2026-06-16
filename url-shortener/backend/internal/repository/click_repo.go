package repository

import (
	"context"
	"database/sql"
	"url-shortener/internal/domain"
)

type ClickRepository interface {
	RecordClick(ctx context.Context, stat *domain.ClickStat) error
	GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error)

	GetLinksByUserID(ctx context.Context, userID int) ([]*domain.Link, error)
	GetClicksByLinkID(ctx context.Context, linkID int, days int) ([]*domain.ClickStat, error)
	GetLinkByCode(ctx context.Context, code string) (*domain.Link, error)
	GetSessionBySessionID(ctx context.Context, sessionID string) (*domain.ClickStat, error)
	UpdateClick(ctx context.Context, stat *domain.ClickStat) error
}

type postgresClickRepo struct {
	db *sql.DB
}

func NewClickRepository(db *sql.DB) ClickRepository {
	return &postgresClickRepo{db: db}
}

func (r *postgresClickRepo) RecordClick(ctx context.Context, stat *domain.ClickStat) error {
	query := `INSERT INTO click_stats 
			  (link_id, ip_address, user_agent, country_code, device_type, browser_name) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		stat.LinkID,
		stat.IPAddress,
		stat.UserAgent,
		stat.CountryCode,
		stat.DeviceType,
		stat.BrowserName)
	return err
}

func (r *postgresClickRepo) GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error) {
	stats := &domain.StatsSummary{
		ByCountry: make(map[string]int64),
		ByDevice:  make(map[string]int64),
		ByBrowser: make(map[string]int64),
	}

	// Общее количество
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM click_stats WHERE link_id = $1`, linkID).
		Scan(&stats.TotalClicks)
	if err != nil {
		return nil, err
	}

	// По странам
	rows, err := r.db.QueryContext(ctx, `SELECT country_code, COUNT(*) FROM click_stats 
										WHERE link_id = $1 AND country_code IS NOT NULL 
										GROUP BY country_code`, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var country string
		var count int64
		if err := rows.Scan(&country, &count); err != nil {
			return nil, err
		}
		stats.ByCountry[country] = count
	}

	// По устройствам
	rows, err = r.db.QueryContext(ctx, `SELECT device_type, COUNT(*) FROM click_stats 
									   WHERE link_id = $1 AND device_type IS NOT NULL 
									   GROUP BY device_type`, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var device string
		var count int64
		if err := rows.Scan(&device, &count); err != nil {
			return nil, err
		}
		stats.ByDevice[device] = count
	}

	// По браузерам
	rows, err = r.db.QueryContext(ctx, `SELECT browser_name, COUNT(*) FROM click_stats 
									   WHERE link_id = $1 AND browser_name IS NOT NULL 
									   GROUP BY browser_name`, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var browser string
		var count int64
		if err := rows.Scan(&browser, &count); err != nil {
			return nil, err
		}
		stats.ByBrowser[browser] = count
	}

	return stats, nil
}

func (r *postgresClickRepo) GetLinksByUserID(ctx context.Context, userID int) ([]*domain.Link, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, original_url, short_code, created_at, expires_at 
         FROM links WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*domain.Link
	for rows.Next() {
		var link domain.Link
		err := rows.Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode,
			&link.CreatedAt, &link.ExpiresAt)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, rows.Err()
}

// GetClicksByLinkID - клики по ссылке за период
func (r *postgresClickRepo) GetClicksByLinkID(ctx context.Context, linkID int, days int) ([]*domain.ClickStat, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, link_id, ip_address, country_code, device_type, browser_name, 
                clicked_at, session_id, page_views, time_on_site, is_bounce
         FROM click_stats 
         WHERE link_id = $1 AND clicked_at >= NOW() - INTERVAL '1 day' * $2
         ORDER BY clicked_at`,
		linkID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []*domain.ClickStat
	for rows.Next() {
		var click domain.ClickStat
		err := rows.Scan(&click.ID, &click.LinkID, &click.IPAddress, &click.CountryCode,
			&click.DeviceType, &click.BrowserName, &click.ClickedAt,
			&click.SessionID, &click.PageViews, &click.TimeOnSite, &click.IsBounce)
		if err != nil {
			return nil, err
		}
		clicks = append(clicks, &click)
	}
	return clicks, rows.Err()
}

// GetLinkByCode - найти ссылку по short_code
func (r *postgresClickRepo) GetLinkByCode(ctx context.Context, code string) (*domain.Link, error) {
	var link domain.Link
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, original_url, short_code, created_at, expires_at 
         FROM links WHERE short_code = $1`, code).
		Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode,
			&link.CreatedAt, &link.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, ErrLinkNotFound
	}
	return &link, err
}

// GetSessionBySessionID - найти сессию
func (r *postgresClickRepo) GetSessionBySessionID(ctx context.Context, sessionID string) (*domain.ClickStat, error) {
	if sessionID == "" {
		return nil, nil
	}

	var click domain.ClickStat
	err := r.db.QueryRowContext(ctx,
		`SELECT id, link_id, ip_address, country_code, device_type, browser_name,
                clicked_at, session_id, page_views, time_on_site, is_bounce
         FROM click_stats WHERE session_id = $1 ORDER BY clicked_at DESC LIMIT 1`,
		sessionID).
		Scan(&click.ID, &click.LinkID, &click.IPAddress, &click.CountryCode,
			&click.DeviceType, &click.BrowserName, &click.ClickedAt,
			&click.SessionID, &click.PageViews, &click.TimeOnSite, &click.IsBounce)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &click, err
}

// UpdateClick - обновить запись клика
func (r *postgresClickRepo) UpdateClick(ctx context.Context, stat *domain.ClickStat) error {
	query := `UPDATE click_stats 
        SET page_views = $1, time_on_site = $2, is_bounce = $3, 
            last_activity_at = $4
        WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		stat.PageViews, stat.TimeOnSite, stat.IsBounce,
		stat.LastActivityAt, stat.ID)

	return err
}
