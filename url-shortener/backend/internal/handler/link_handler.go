package handler

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	sessionRepo      repository.SessionRepository
	baseURL          string
}

// 🔥 Конструктор с sessionRepo
func NewLinkHandler(linkService service.LinkService, analyticsService service.AnalyticsService, sessionRepo repository.SessionRepository, baseURL string) *LinkHandler {
	return &LinkHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
		sessionRepo:      sessionRepo, // ← Теперь работает
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

	// Асинхронно записываем аналитику (не блокируем ответ)
	go h.recordSession(r.Context(), link.ID, r)

	WriteJSON(w, http.StatusOK, map[string]string{
		"original_url": link.OriginalURL,
	})
}

// 🔥 ЕДИНЫЙ метод Redirect (удали дубликат!)
// Путь: /{code}
func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" || code == "/" || strings.HasPrefix(code, "api/") {
		http.NotFound(w, r)
		return
	}

	// 1. Ищем ссылку в БД
	link, err := h.linkService.GetLinkByCode(r.Context(), code)
	if err != nil {
		log.Printf("⚠️ Link not found for code: %s", code)
		http.NotFound(w, r)
		return
	}

	log.Printf("🚀 Redirecting code: %s -> %s", code, link.OriginalURL)

	// 2. Записываем аналитику АСИНХРОННО
	// 🔥 ВАЖНО: используем context.Background(), иначе запись прервется после редиректа
	go h.recordSession(context.Background(), link.ID, r)

	// 3. Делаем редирект
	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

// 🔥 recordSession — запись аналитики при переходе по ссылке
func (h *LinkHandler) recordSession(ctx context.Context, linkID int, r *http.Request) {
	// 🔥 Создаем НОВЫЙ контекст с таймаутом, независимый от запроса
	recordCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("📝 [Async] Starting recordSession for linkID=%d...", linkID)

	// 1. Получаем IP
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}
	ipHash := hashIP(strings.TrimSpace(ip))

	// 2. UserAgent
	ua := r.UserAgent()
	deviceType := "desktop"
	if strings.Contains(strings.ToLower(ua), "mobile") ||
		strings.Contains(strings.ToLower(ua), "android") ||
		strings.Contains(strings.ToLower(ua), "iphone") {
		deviceType = "mobile"
	}
	browser := parseBrowser(ua)
	country := "RU"

	// 3. SessionID
	sessionRaw := ipHash + ua
	hash := sha256.Sum256([]byte(sessionRaw))
	sessionID := fmt.Sprintf("%x", hash)

	log.Printf("📝 [Async] Trying to INSERT into analytics_sessions...")

	// 4. Пишем в БД с новым контекстом
	err := h.sessionRepo.RecordSession(recordCtx, linkID, sessionID, ipHash, country, deviceType, browser)
	if err != nil {
		log.Printf("❌ [Async] DB Error: %v", err)
	} else {
		log.Println("✅ [Async] Session recorded successfully!")
	}
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "short code is required")
		return
	}

	err := h.linkService.DeleteLink(r.Context(), userID, code)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			WriteError(w, http.StatusNotFound, "link not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
