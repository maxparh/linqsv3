package domain

import "time"

type ClickStat struct {
	ID          int       `json:"id"`
	LinkID      int       `json:"link_id"`
	IPAddress   string    `json:"-"`
	UserAgent   string    `json:"user_agent"`
	CountryCode string    `json:"country_code"`
	DeviceType  string    `json:"device_type"`
	BrowserName string    `json:"browser_name"`
	ClickedAt   time.Time `json:"clicked_at"`

	SessionID      string    `json:"session_id"`
	PageViews      int       `json:"page_views"`
	TimeOnSite     int       `json:"time_on_site"`
	IsBounce       bool      `json:"is_bounce"`
	LastActivityAt time.Time `json:"last_activity_at"`
}

type StatsSummary struct {
	TotalClicks int64            `json:"total_clicks"`
	ByCountry   map[string]int64 `json:"by_country"`
	ByDevice    map[string]int64 `json:"by_device"`
	ByBrowser   map[string]int64 `json:"by_browser"`
}

type ClickByDay struct {
	Day    time.Time `json:"day"`
	Clicks int64     `json:"clicks"`
}

type LinkStats struct {
	LinkID      int              `json:"link_id"`
	ShortCode   string           `json:"short_code"`
	OriginalURL string           `json:"original_url"`
	CreatedAt   time.Time        `json:"created_at"`
	TotalClicks int64            `json:"total_clicks"`
	ClicksByDay []ClickByDay     `json:"clicks_by_day"`
	ByCountry   map[string]int64 `json:"by_country"`
}
