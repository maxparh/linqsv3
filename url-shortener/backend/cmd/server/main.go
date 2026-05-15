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

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	linkService := service.NewLinkService(linkRepo)
	analyticsService := service.NewAnalyticsService(clickRepo)

	// Инициализация хендлеров
	authHandler := handler.NewAuthHandler(authService)
	linkHandler := handler.NewLinkHandler(linkService, analyticsService, cfg.FrontendURL) // Или отдельный домен для ссылок

	// Настройка роутера
	mux := http.NewServeMux()

	// Публичные роуты (с указанием метода)
	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	// 🔥 Middleware авторизации
	authMiddleware := middleware.NewAuthMiddleware(authService)

	mux.HandleFunc("GET /api/links", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(linkHandler.GetUserLinks)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/links", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(linkHandler.CreateLink)).ServeHTTP(w, r)
	})

	mux.HandleFunc("DELETE /api/links/{code}", func(w http.ResponseWriter, r *http.Request) {
		authMiddleware.RequireAuth(http.HandlerFunc(linkHandler.DeleteLink)).ServeHTTP(w, r)
	})

	// Публичный эндпоинт для резолва
	mux.HandleFunc("GET /api/resolve/{code}", linkHandler.ResolveLink)

	// 🔥 Редиректы — ОБЯЗАТЕЛЬНО в конце, ловит всё остальное
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
