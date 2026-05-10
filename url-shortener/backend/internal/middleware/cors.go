package middleware

import (
	"net/http"
)

type CORSMiddleware struct {
	allowedOrigin string
}

func NewCORSMiddleware(allowedOrigin string) *CORSMiddleware {
	return &CORSMiddleware{allowedOrigin: allowedOrigin}
}

func (m *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", m.allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если это preflight запрос (OPTIONS), сразу отвечаем OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
