package repository

import (
	"context"
	"database/sql"
	"url-shortener/internal/domain"
)

type ClickRepository interface {
	RecordClick(ctx context.Context, stat *domain.ClickStat) error
	GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error)
}

type postgresClickRepo struct {
	db *sql.DB
}

func NewClickRepository(db *sql.DB) ClickRepository {
	return &postgresClickRepo{db: db}
}

func (r *postgresClickRepo) RecordClick(ctx context.Context, stat *domain.ClickStat) error {
	// Используем INSERT ... ON CONFLICT DO NOTHING если нужно избегать дублей по IP+Link за короткое время,
	// но для простоты пишем всё подряд. В продакшене можно агрегировать в Redis перед записью в БД.
	query := `INSERT INTO click_stats (link_id, ip_address, user_agent, country_code, device_type, browser_name) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		stat.LinkID,
		stat.IPAddress,
		stat.UserAgent,
		stat.Country,
		stat.DeviceType,
		stat.Browser)
	return err
}

func (r *postgresClickRepo) GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error) {
	stats := &domain.StatsSummary{
		ByCountry: make(map[string]int64),
		ByDevice:  make(map[string]int64),
	}

	// Общее количество
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM click_stats WHERE link_id = $1`, linkID).
		Scan(&stats.TotalClicks)
	if err != nil {
		return nil, err
	}

	// Группировка по странам
	rows, err := r.db.QueryContext(ctx, `SELECT country_code, COUNT(*) FROM click_stats 
										WHERE link_id = $1 GROUP BY country_code`, linkID)
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

	// Группировка по устройствам
	rows, err = r.db.QueryContext(ctx, `SELECT device_type, COUNT(*) FROM click_stats 
									   WHERE link_id = $1 GROUP BY device_type`, linkID)
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

	return stats, nil
}
