package domain

import "time"

// AnalyticsSession представляет сессию пользователя
type AnalyticsSession struct {
	ID              int       `json:"id"`
	SessionID       string    `json:"session_id"`
	LinkID          int       `json:"link_id"`
	UserIPHash      string    `json:"user_ip_hash"`
	StartedAt       time.Time `json:"started_at"`
	LastActivityAt  time.Time `json:"last_activity_at"`
	PageViews       int       `json:"page_views"`
	TotalTimeOnSite int       `json:"total_time_on_site"` // в секундах
	IsBounce        bool      `json:"is_bounce"`
	CountryCode     string    `json:"country_code"`
	DeviceType      string    `json:"device_type"`
	BrowserName     string    `json:"browser_name"`
}

// CreateSessionRequest - данные для создания сессии
type CreateSessionRequest struct {
	SessionID   string `json:"session_id"`
	LinkID      int    `json:"link_id"`
	UserIPHash  string `json:"user_ip_hash"`
	CountryCode string `json:"country_code"`
	DeviceType  string `json:"device_type"`
	BrowserName string `json:"browser_name"`
}

// UpdateSessionRequest - обновление сессии
type UpdateSessionRequest struct {
	SessionID       string `json:"session_id"`
	PageViews       int    `json:"page_views"`
	TotalTimeOnSite int    `json:"total_time_on_site"`
	IsBounce        bool   `json:"is_bounce"`
}
