package service

import (
	"context"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
)

type AnalyticsService interface {
	RecordClick(ctx context.Context, stat *domain.ClickStat) error
	GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error)
}

type analyticsService struct {
	clickRepo repository.ClickRepository
}

func NewAnalyticsService(clickRepo repository.ClickRepository) AnalyticsService {
	return &analyticsService{clickRepo: clickRepo}
}

// RecordClick передаёт данные клика в репозиторий
func (s *analyticsService) RecordClick(ctx context.Context, stat *domain.ClickStat) error {
	return s.clickRepo.RecordClick(ctx, stat)
}

// GetStats возвращает агрегированную статистику по ссылке
func (s *analyticsService) GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error) {
	return s.clickRepo.GetStats(ctx, linkID)
}