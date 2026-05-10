package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"url-shortener/internal/domain"
	"url-shortener/internal/service"
	"url-shortener/internal/middleware"
)

type LinkHandler struct {
	linkService      service.LinkService
	analyticsService service.AnalyticsService
	baseURL          string // Например, https://short.ly
}

func NewLinkHandler(linkService service.LinkService, analyticsService service.AnalyticsService, baseURL string) *LinkHandler {
	return &LinkHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
		baseURL:          baseURL,
	}
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Получаем UserID из контекста (добавим middleware позже)
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

// Redirect обрабатывает переход по короткой ссылке
func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Извлекаем код из URL. Ожидаем путь вида /{code}
	code := strings.TrimPrefix(r.URL.Path, "/")

	// Если код пустой или это корень, отдаем 404 или главную
	if code == "" || code == "/" {
		http.NotFound(w, r)
		return
	}

	link, err := h.linkService.GetLinkByCode(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Асинхронно записываем клик, чтобы не задерживать редирект
	go h.recordClick(r.Context(), link.ID, r)

	// Выполняем редирект
	http.Redirect(w, r, link.OriginalURL, http.StatusFound) // 302 Found
}

func (h *LinkHandler) recordClick(ctx context.Context, linkID int, r *http.Request) {
	// Простая эвристика для определения устройства
	deviceType := "desktop"
	if strings.Contains(strings.ToLower(r.UserAgent()), "mobile") {
		deviceType = "mobile"
	}

	// Парсинг IP (учитывая прокси)
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}

	stat := &domain.ClickStat{
		LinkID:     linkID,
		IPAddress:  ip,
		UserAgent:  r.UserAgent(),
		Country:    "US", // Заглушка, в реальном проекте нужен GeoIP сервис
		DeviceType: deviceType,
		Browser:    "Unknown", // Заглушка
	}

	h.analyticsService.RecordClick(context.Background(), stat)
}
