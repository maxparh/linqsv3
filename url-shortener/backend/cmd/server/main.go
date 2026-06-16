package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func main() {
	cfg := config.Load()

	// Инициализация БД
	db, err := repository.NewPostgresDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db.DB)
	linkRepo := repository.NewLinkRepository(db.DB)
	clickRepo := repository.NewClickRepository(db.DB)
	sessionRepo := repository.NewSessionRepository(db.DB)

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	linkService := service.NewLinkService(linkRepo)

	// Передаём оба репозитория
	analyticsService := service.NewAnalyticsService(clickRepo, sessionRepo)

	// Инициализация хендлеров
	authHandler := handler.NewAuthHandler(authService)

	// Передаём sessionRepo
	linkHandler := handler.NewLinkHandler(linkService, analyticsService, sessionRepo, cfg.FrontendURL)

	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Настройка роутера
	mux := http.NewServeMux()

	// Публичные роуты
	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	cookieHandler := handler.NewCookieConsentHandler()
	mux.HandleFunc("POST /api/cookies/consent", cookieHandler.SaveConsent)

	// Middleware авторизации
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Защищенные роуты для ссылок
	mux.HandleFunc("GET /api/links", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(linkHandler.GetUserLinks)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/links", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(linkHandler.CreateLink)).ServeHTTP(w, r)
	})

	mux.HandleFunc("DELETE /api/links/{code}", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(linkHandler.DeleteLink)).ServeHTTP(w, r)
	})

	// Защищённые роуты для аналитики
	mux.HandleFunc("GET /api/analytics/overview", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(analyticsHandler.GetOverview)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/analytics/clicks-over-time", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(analyticsHandler.GetClicksOverTime)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/analytics/top-locations", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(analyticsHandler.GetTopLocations)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/analytics/devices", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(analyticsHandler.GetDeviceStats)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/analytics/bounce-rate", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(analyticsHandler.GetBounceRate)).ServeHTTP(w, r)
	})

	// Публичные эндпоинты
	mux.HandleFunc("GET /api/resolve/{code}", linkHandler.ResolveLink)
	mux.HandleFunc("POST /api/analytics/track", analyticsHandler.TrackSession)

	// Редиректы — в конце
	mux.HandleFunc("/", linkHandler.Redirect)

	// CORS Middleware
	corsMiddleware := middleware.NewCORSMiddleware(cfg.FrontendURL)
	handlerWithCORS := corsMiddleware.Handler(mux)

	// Сервер
	srv := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      handlerWithCORS,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful Shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
