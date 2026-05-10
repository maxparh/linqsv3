package service

import (
	"context"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
)

type AnalyticsService interface {
	RecordClick(ctx context.Context, stat *domain.ClickStat)
	GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error)
}

type analyticsService struct {
	clickRepo repository.ClickRepository
}

func NewAnalyticsService(clickRepo repository.ClickRepository) AnalyticsService {
	return &analyticsService{clickRepo: clickRepo}
}

func (s *analyticsService) RecordClick(ctx context.Context, stat *domain.ClickStat) {
	// В продакшене здесь можно использовать канал (channel) для асинхронной записи,
	// чтобы не блокировать редирект. Но для простоты пишем синхронно.
	// Если БД упадет, мы просто потеряем клик, но пользователь увидит редирект.
	_ = s.clickRepo.RecordClick(ctx, stat)
}

func (s *analyticsService) GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error) {
	return s.clickRepo.GetStats(ctx, linkID)
}
