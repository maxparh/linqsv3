package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"url-shortener/internal/domain"
	"url-shortener/internal/middleware"
	"url-shortener/internal/service"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetOverview - GET /api/analytics/overview?days=7
func (h *AnalyticsHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		days, _ = strconv.Atoi(d)
		if days <= 0 {
			days = 7
		}
	}

	overview, err := h.analyticsService.GetOverview(r.Context(), userID, days)
	if err != nil {
		log.Printf("❌ Analytics Error: %v", err)
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("✅ Analytics Overview: %+v", overview) // ← Добавь это
	WriteJSON(w, http.StatusOK, overview)
}

// GetClicksOverTime - GET /api/analytics/clicks-over-time?days=7
func (h *AnalyticsHandler) GetClicksOverTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		days, _ = strconv.Atoi(d)
		if days <= 0 {
			days = 7
		}
	}

	data, err := h.analyticsService.GetClicksOverTime(r.Context(), userID, days)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, data)
}

// GetTopLocations - GET /api/analytics/top-locations?days=7&limit=5
func (h *AnalyticsHandler) GetTopLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		days, _ = strconv.Atoi(d)
		if days <= 0 {
			days = 7
		}
	}

	limit := 5
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
		if limit <= 0 {
			limit = 5
		}
	}

	locations, err := h.analyticsService.GetTopLocations(r.Context(), userID, days, limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, locations)
}

// GetDeviceStats - GET /api/analytics/devices?days=7
func (h *AnalyticsHandler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		days, _ = strconv.Atoi(d)
		if days <= 0 {
			days = 7
		}
	}

	devices, err := h.analyticsService.GetDeviceStats(r.Context(), userID, days)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, devices)
}

// GetBounceRate - GET /api/analytics/bounce-rate?days=7
func (h *AnalyticsHandler) GetBounceRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		days, _ = strconv.Atoi(d)
		if days <= 0 {
			days = 7
		}
	}

	bounceRate, err := h.analyticsService.GetBounceRate(r.Context(), userID, days)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, bounceRate)
}

// TrackSession - POST /api/analytics/track (публичный)
func (h *AnalyticsHandler) TrackSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		SessionID   string `json:"session_id"`
		LinkCode    string `json:"link_code"`
		UserIPHash  string `json:"user_ip_hash"`
		CountryCode string `json:"country_code"`
		DeviceType  string `json:"device_type"`
		BrowserName string `json:"browser_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	trackReq := &domain.TrackSessionRequest{
		SessionID:   req.SessionID,
		LinkCode:    req.LinkCode,
		UserIPHash:  req.UserIPHash,
		CountryCode: req.CountryCode,
		DeviceType:  req.DeviceType,
		BrowserName: req.BrowserName,
	}

	if err := h.analyticsService.TrackSession(r.Context(), trackReq); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
