package domain

import "time"

type Link struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // nil если нет срока
}

// CreateLinkRequest приходит от пользователя
type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"`
	// ShortCode можно оставить пустым, тогда сгенерируем сами
	// Или пользователь может предложить свой (нужна проверка на уникальность)
	ShortCode string `json:"short_code,omitempty"`
}

// LinkResponse отдается пользователю после создания
type LinkResponse struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	FullURL     string `json:"full_url"` // Например: https://mysite.com/aB3dE
}
