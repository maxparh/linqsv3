package domain

type AnalyticsOverview struct {
	TotalClicks   int64   `json:"total_clicks"`
	UniqueClicks  int64   `json:"unique_clicks"`
	BounceRate    float64 `json:"bounce_rate"`
	AvgTimeOnSite float64 `json:"avg_time_on_site"`
}

type ClicksOverTime struct {
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

type LocationStat struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Percent     float64 `json:"percent"`
}

type DeviceStat struct {
	Name    string  `json:"name"`
	Value   int     `json:"value"`
	Percent float64 `json:"percent"`
	Color   string  `json:"color"`
}

type BounceRateStat struct {
	LinkURL       string  `json:"link_url"`
	BounceRate    float64 `json:"bounce_rate"`
	TotalSessions int     `json:"total_sessions"`
}

type TrackSessionRequest struct {
	SessionID   string `json:"session_id"`
	LinkCode    string `json:"link_code"`
	UserIPHash  string `json:"user_ip_hash"`
	CountryCode string `json:"country_code"`
	DeviceType  string `json:"device_type"`
	BrowserName string `json:"browser_name"`
}
