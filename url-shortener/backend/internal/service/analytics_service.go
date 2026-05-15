package service

import (
	"context"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
)

type AnalyticsService interface {
	RecordClick(ctx context.Context, stat *domain.ClickStat) error
	GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error)

	GetOverview(ctx context.Context, userID int, days int) (*domain.AnalyticsOverview, error)
	GetClicksOverTime(ctx context.Context, userID int, days int) (*domain.ClicksOverTime, error)
	GetTopLocations(ctx context.Context, userID int, days int, limit int) ([]*domain.LocationStat, error)
	GetDeviceStats(ctx context.Context, userID int, days int) ([]*domain.DeviceStat, error)
	GetBounceRate(ctx context.Context, userID int, days int) ([]*domain.BounceRateStat, error)
	TrackSession(ctx context.Context, req *domain.TrackSessionRequest) error
}

type analyticsService struct {
	clickRepo   repository.ClickRepository
	sessionRepo repository.SessionRepository
}

func NewAnalyticsService(clickRepo repository.ClickRepository, sessionRepo repository.SessionRepository) AnalyticsService {
	return &analyticsService{
		clickRepo:   clickRepo,
		sessionRepo: sessionRepo,
	}
}

// RecordClick передаёт данные клика в репозиторий
func (s *analyticsService) RecordClick(ctx context.Context, stat *domain.ClickStat) error {
	return s.clickRepo.RecordClick(ctx, stat)
}

// GetStats возвращает агрегированную статистику по ссылке
func (s *analyticsService) GetStats(ctx context.Context, linkID int) (*domain.StatsSummary, error) {
	return s.clickRepo.GetStats(ctx, linkID)
}

func (s *analyticsService) GetOverview(ctx context.Context, userID int, days int) (*domain.AnalyticsOverview, error) {
	return s.sessionRepo.GetOverviewStats(ctx, userID, days)
}

func (s *analyticsService) GetClicksOverTime(ctx context.Context, userID int, days int) (*domain.ClicksOverTime, error) {
	return s.sessionRepo.GetClicksOverTime(ctx, userID, days)
}

func (s *analyticsService) GetTopLocations(ctx context.Context, userID int, days int, limit int) ([]*domain.LocationStat, error) {
	return s.sessionRepo.GetTopLocations(ctx, userID, days, limit)
}

func (s *analyticsService) GetDeviceStats(ctx context.Context, userID int, days int) ([]*domain.DeviceStat, error) {
	return s.sessionRepo.GetDeviceStats(ctx, userID, days)
}

func (s *analyticsService) GetBounceRate(ctx context.Context, userID int, days int) ([]*domain.BounceRateStat, error) {
	ov, err := s.GetOverview(ctx, userID, days)
	if err != nil {
		return nil, err
	}
	return []*domain.BounceRateStat{{
		LinkURL:       "Общая статистика",
		BounceRate:    ov.BounceRate,
		TotalSessions: int(ov.UniqueClicks),
	}}, nil
}

// TrackSession - создание/обновление сессии
func (s *analyticsService) TrackSession(ctx context.Context, req *domain.TrackSessionRequest) error {
	// Находим link_id по short_code
	link, err := s.clickRepo.GetLinkByCode(ctx, req.LinkCode)
	if err != nil {
		return err
	}

	// Проверяем существующую сессию
	existing, err := s.clickRepo.GetSessionBySessionID(ctx, req.SessionID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Новая сессия
		click := &domain.ClickStat{
			LinkID:      link.ID,
			SessionID:   req.SessionID,
			IPAddress:   req.UserIPHash,
			CountryCode: req.CountryCode,
			DeviceType:  req.DeviceType,
			BrowserName: req.BrowserName,
			ClickedAt:   time.Now(),
			IsBounce:    true,
			PageViews:   1,
			TimeOnSite:  0,
		}
		return s.clickRepo.RecordClick(ctx, click)
	}

	// Обновляем существующую
	now := time.Now()
	timeOnSite := int(now.Sub(existing.ClickedAt).Seconds())

	existing.PageViews++
	existing.TimeOnSite += timeOnSite
	existing.LastActivityAt = now
	existing.IsBounce = false

	return s.clickRepo.UpdateClick(ctx, existing)
}

func (s *analyticsService) RecordSession(ctx context.Context, linkID int, sessionID string, ipHash string, countryCode string, deviceType string, browserName string) error {
	return s.sessionRepo.RecordSession(ctx, linkID, sessionID, ipHash, countryCode, deviceType, browserName)
}
