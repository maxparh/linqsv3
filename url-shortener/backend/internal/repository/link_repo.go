package repository

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/internal/domain"
)

var ErrLinkNotFound = errors.New("link not found")
var ErrCodeExists = errors.New("short code already exists")

type LinkRepository interface {
	Create(ctx context.Context, link *domain.Link) error
	GetByCode(ctx context.Context, code string) (*domain.Link, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Link, error)
}

type postgresLinkRepo struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) LinkRepository {
	return &postgresLinkRepo{db: db}
}

func (r *postgresLinkRepo) Create(ctx context.Context, link *domain.Link) error {
	query := `INSERT INTO links (user_id, original_url, short_code, expires_at) 
			  VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, link.UserID, link.OriginalURL, link.ShortCode, link.ExpiresAt).
		Scan(&link.ID, &link.CreatedAt)

	// Проверка на уникальность кода (если вдруг коллизия)
	if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"links_short_code_key\"" {
		return ErrCodeExists
	}
	return err
}

func (r *postgresLinkRepo) GetByCode(ctx context.Context, code string) (*domain.Link, error) {
	var link domain.Link
	query := `SELECT id, user_id, original_url, short_code, created_at, expires_at 
			  FROM links WHERE short_code = $1`
	err := r.db.QueryRowContext(ctx, query, code).
		Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.CreatedAt, &link.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, ErrLinkNotFound
	}
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *postgresLinkRepo) GetByUserID(ctx context.Context, userID int) ([]*domain.Link, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, original_url, short_code, created_at, expires_at 
		FROM links WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*domain.Link
	for rows.Next() {
		var link domain.Link
		if err := rows.Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.CreatedAt, &link.ExpiresAt); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, rows.Err()
}
