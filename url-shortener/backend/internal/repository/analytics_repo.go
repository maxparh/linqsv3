package repository

import (
	"context"
	"database/sql"
	"url-shortener/internal/domain"
)

type AnalyticsRepository interface {
	GetUserStatsSummary(ctx context.Context, userID int) (*domain.StatsSummary, error)
}

type postgresAnalyticsRepo struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) AnalyticsRepository {
	return &postgresAnalyticsRepo{db: db}
}

func (r *postgresAnalyticsRepo) GetUserStatsSummary(ctx context.Context, userID int) (*domain.StatsSummary, error) {
	summary := &domain.StatsSummary{
		ByCountry: make(map[string]int64),
		ByDevice:  make(map[string]int64),
		ByBrowser: make(map[string]int64),
	}

	// Всего кликов по ссылкам пользователя
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM click_stats cs
		JOIN links l ON cs.link_id = l.id
		WHERE l.user_id = $1
	`, userID).Scan(&summary.TotalClicks)
	if err != nil {
		return nil, err
	}

	// По странам
	rows, err := r.db.QueryContext(ctx, `
		SELECT country_code, COUNT(*) as cnt
		FROM click_stats cs
		JOIN links l ON cs.link_id = l.id
		WHERE l.user_id = $1 AND country_code IS NOT NULL AND country_code != ''
		GROUP BY country_code
		ORDER BY cnt DESC
		LIMIT 10
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var country string
		var cnt int64
		if err := rows.Scan(&country, &cnt); err != nil {
			return nil, err
		}
		summary.ByCountry[country] = cnt
	}

	// По устройствам
	rows, err = r.db.QueryContext(ctx, `
		SELECT device_type, COUNT(*) as cnt
		FROM click_stats cs
		JOIN links l ON cs.link_id = l.id
		WHERE l.user_id = $1 AND device_type IS NOT NULL AND device_type != ''
		GROUP BY device_type
		ORDER BY cnt DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var device string
		var cnt int64
		if err := rows.Scan(&device, &cnt); err != nil {
			return nil, err
		}
		summary.ByDevice[device] = cnt
	}

	// По браузерам
	rows, err = r.db.QueryContext(ctx, `
		SELECT browser_name, COUNT(*) as cnt
		FROM click_stats cs
		JOIN links l ON cs.link_id = l.id
		WHERE l.user_id = $1 AND browser_name IS NOT NULL AND browser_name != ''
		GROUP BY browser_name
		ORDER BY cnt DESC
		LIMIT 10
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var browser string
		var cnt int64
		if err := rows.Scan(&browser, &cnt); err != nil {
			return nil, err
		}
		summary.ByBrowser[browser] = cnt
	}

	return summary, nil
}