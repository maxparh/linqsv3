package domain

import "time"

type ClickStat struct {
	ID         int       `json:"id"`
	LinkID     int       `json:"link_id"`
	IPAddress  string    `json:"ip_address"` // В БД INET, в Go string для простоты сериализации
	UserAgent  string    `json:"user_agent"`
	Country    string    `json:"country"`
	DeviceType string    `json:"device_type"`
	Browser    string    `json:"browser"`
	ClickedAt  time.Time `json:"clicked_at"`
}

// StatsSummary - агрегированные данные для дашборда
type StatsSummary struct {
	TotalClicks int64            `json:"total_clicks"`
	ByCountry   map[string]int64 `json:"by_country"`
	ByDevice    map[string]int64 `json:"by_device"`
}
