package handler

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

type LinkHandler struct {
	linkService      service.LinkService
	analyticsService service.AnalyticsService
	baseURL          string
}

func NewLinkHandler(linkService service.LinkService, analyticsService service.AnalyticsService, baseURL string) *LinkHandler {
	return &LinkHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
		baseURL:          baseURL,
	}
}

// hashIP хэширует IP через SHA256 (64 символа)
func hashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return fmt.Sprintf("%x", hash)
}

// parseBrowser определяет браузер из User-Agent
func parseBrowser(ua string) string {
	uaLower := strings.ToLower(ua)
	switch {
	case strings.Contains(uaLower, "edg"):
		return "Edge"
	case strings.Contains(uaLower, "chrome"):
		return "Chrome"
	case strings.Contains(uaLower, "firefox"):
		return "Firefox"
	case strings.Contains(uaLower, "safari"):
		return "Safari"
	case strings.Contains(uaLower, "opera") || strings.Contains(uaLower, "opr"):
		return "Opera"
	default:
		return "Unknown"
	}
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req domain.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.linkService.CreateLink(r.Context(), userID, &req, h.baseURL)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (h *LinkHandler) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	links, err := h.linkService.GetUserLinks(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, links)
}

// ResolveLink — API-эндпоинт для получения original_url (для фронтенда)
// Путь: /api/resolve/{code}
func (h *LinkHandler) ResolveLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/api/resolve/")
	if code == "" || strings.Contains(code, "/") {
		WriteError(w, http.StatusBadRequest, "invalid short code")
		return
	}

	link, err := h.linkService.GetLinkByCode(r.Context(), code)
	if err != nil {
		WriteError(w, http.StatusNotFound, "link not found")
		return
	}

	// Асинхронно записываем клик (не блокируем ответ)
	go h.recordClick(r.Context(), link.ID, r)

	WriteJSON(w, http.StatusOK, map[string]string{
		"original_url": link.OriginalURL,
	})
}

// Redirect — прямой редирект (для старых ссылок или если фронт не используется)
// Путь: /{code}
func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" || code == "/" || strings.HasPrefix(code, "api/") {
		http.NotFound(w, r)
		return
	}

	link, err := h.linkService.GetLinkByCode(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	go h.recordClick(r.Context(), link.ID, r)
	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

// recordClick — запись статистики клика
func (h *LinkHandler) recordClick(ctx context.Context, linkID int, r *http.Request) {
	// Устройство
	deviceType := "desktop"
	ua := r.UserAgent()
	uaLower := strings.ToLower(ua)
	if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "android") || strings.Contains(uaLower, "iphone") {
		deviceType = "mobile"
	}

	// IP (с учётом прокси)
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}
	ipHash := hashIP(strings.TrimSpace(ip))

	// Браузер
	browser := parseBrowser(ua)

	// Создаём объект клика с полями, соответствующими БД
	stat := &domain.ClickStat{
		LinkID:      linkID,
		IPAddress:   ipHash, // SHA256 хэш (64 символа)
		UserAgent:   ua,
		CountryCode: "RU", // ← Поле как в БД: country_code
		DeviceType:  deviceType,
		BrowserName: browser,     // ← Поле как в БД: browser_name
		ClickedAt:   time.Time{}, // Zero value → БД подставит NOW()
	}

	// Записываем клик (игнорируем ошибку, чтобы не ломать редирект)
	_ = h.analyticsService.RecordClick(context.Background(), stat)
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	// Метод уже проверен на уровне роутера, но для надёжности:
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Достаём user_id из контекста (как у тебя уже сделано)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// 🔥 Go 1.22+: достаём параметр из пути через r.PathValue()
	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "short code is required")
		return
	}

	// Вызываем сервис
	err := h.linkService.DeleteLink(r.Context(), userID, code)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			WriteError(w, http.StatusNotFound, "link not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Успех — 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
