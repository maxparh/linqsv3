package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
)

const shortCodeLength = 6

type LinkService interface {
	CreateLink(ctx context.Context, userID int, req *domain.CreateLinkRequest, baseURL string) (*domain.LinkResponse, error)
	GetLinkByCode(ctx context.Context, code string) (*domain.Link, error)
	GetUserLinks(ctx context.Context, userID int) ([]*domain.Link, error)
	DeleteLink(ctx context.Context, userID int, code string) error
}

type linkService struct {
	linkRepo repository.LinkRepository
}

func NewLinkService(linkRepo repository.LinkRepository) LinkService {
	return &linkService{linkRepo: linkRepo}
}

func (s *linkService) CreateLink(ctx context.Context, userID int, req *domain.CreateLinkRequest, baseURL string) (*domain.LinkResponse, error) {
	// Валидация URL
	if _, err := url.ParseRequestURI(req.OriginalURL); err != nil {
		return nil, errors.New("invalid original URL")
	}

	shortCode := req.ShortCode
	if shortCode == "" {
		// Генерация случайного кода
		var err error
		shortCode, err = s.generateShortCode()
		if err != nil {
			return nil, err
		}
	} else if len(shortCode) > 10 {
		return nil, errors.New("short code too long")
	}

	link := &domain.Link{
		UserID:      userID,
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
	}

	// Попытка создать ссылку. Если код занят, пробуем снова (простая стратегия)
	for i := 0; i < 3; i++ {
		err := s.linkRepo.Create(ctx, link)
		if err == nil {
			break
		}
		if errors.Is(err, repository.ErrCodeExists) {
			// Генерируем новый код и пробуем снова
			shortCode, _ = s.generateShortCode()
			link.ShortCode = shortCode
			continue
		}
		return nil, err
	}

	fullURL := fmt.Sprintf("%s/%s", baseURL, link.ShortCode)
	return &domain.LinkResponse{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		FullURL:     fullURL,
	}, nil
}

func (s *linkService) generateShortCode() (string, error) {
	b := make([]byte, shortCodeLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Используем base64url encoding без паддинга, чтобы получить безопасные для URL символы
	code := base64.RawURLEncoding.EncodeToString(b)
	// Обрезаем до нужной длины, если base64 выдал больше (маловероятно, но на всякий случай)
	if len(code) > shortCodeLength {
		code = code[:shortCodeLength]
	}
	return code, nil
}

func (s *linkService) GetLinkByCode(ctx context.Context, code string) (*domain.Link, error) {
	return s.linkRepo.GetByCode(ctx, code)
}

func (s *linkService) GetUserLinks(ctx context.Context, userID int) ([]*domain.Link, error) {
	return s.linkRepo.GetByUserID(ctx, userID)
}

func (s *linkService) DeleteLink(ctx context.Context, userID int, code string) error {
	return s.linkRepo.DeleteByCode(ctx, userID, code)
}
