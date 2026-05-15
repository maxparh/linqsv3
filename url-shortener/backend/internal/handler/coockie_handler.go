package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type CookieConsentHandler struct{}

type ConsentRequest struct {
	Accepted    bool            `json:"accepted"`
	Preferences map[string]bool `json:"preferences"`
}

func NewCookieConsentHandler() *CookieConsentHandler {
	return &CookieConsentHandler{}
}

func (h *CookieConsentHandler) SaveConsent(w http.ResponseWriter, r *http.Request) {
	var req ConsentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Создаём cookie с согласием
	consentCookie := &http.Cookie{
		Name:     "cookie_consent",
		Value:    encodeConsent(req),
		Path:     "/",
		Expires:  time.Now().AddDate(1, 0, 0), // 1 год
		HttpOnly: false,                       // false, чтобы фронт мог прочитать при необходимости
		Secure:   r.TLS != nil,                // true для HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, consentCookie)

	// Опционально: сохраняем в БД для аудита
	// userID, _ := r.Context().Value("userID").(string)
	// if userID != "" {
	//     h.consentService.SaveToDB(r.Context(), userID, req)
	// }

	w.WriteHeader(http.StatusNoContent)
}

func encodeConsent(req ConsentRequest) string {
	// Простая кодировка: "accepted:1|analytics:1|marketing:0"
	data, _ := json.Marshal(req)
	return string(data)
}
