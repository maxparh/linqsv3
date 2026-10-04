package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
	"url-shortener/internal/domain"
)

var ErrInvalidProfile = errors.New("invalid profile")
var profilePhonePattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

func validateProfile(req domain.UpdateProfileRequest) (domain.UpdateProfileRequest, error) {
	if req.Avatar != nil && *req.Avatar != "" {
		if err := validateAvatar(*req.Avatar); err != nil {
			return req, fmt.Errorf("%w: invalid avatar", ErrInvalidProfile)
		}
	}
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	for _, name := range []string{req.FirstName, req.LastName} {
		if name == "" || utf8.RuneCountInString(name) > 100 || strings.IndexFunc(name, unicode.IsControl) >= 0 {
			return req, fmt.Errorf("%w: names must contain 1 to 100 characters", ErrInvalidProfile)
		}
	}
	address, err := mail.ParseAddress(req.Email)
	if err != nil || address.Address != req.Email || len(req.Email) > 255 {
		return req, fmt.Errorf("%w: invalid email", ErrInvalidProfile)
	}
	if req.Phone != "" {
		for _, char := range req.Phone {
			if !strings.ContainsRune("+0123456789 ()-.", char) {
				return req, fmt.Errorf("%w: invalid phone", ErrInvalidProfile)
			}
		}
		req.Phone = normalizePhoneToE164(req.Phone)
		if !profilePhonePattern.MatchString(req.Phone) {
			return req, fmt.Errorf("%w: invalid phone", ErrInvalidProfile)
		}
	}
	return req, nil
}

func (s *authService) GetProfile(ctx context.Context, userID int) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func validateAvatar(value string) error {
	header, encoded, ok := strings.Cut(value, ",")
	if !ok || (header != "data:image/jpeg;base64" && header != "data:image/png;base64") || len(encoded) > base64.StdEncoding.EncodedLen(256*1024) {
		return errors.New("invalid image data")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(data) > 256*1024 {
		return errors.New("invalid image size")
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width < 1 || config.Height < 1 || config.Width > 1024 || config.Height > 1024 || header != "data:image/"+format+";base64" {
		return errors.New("invalid image format or dimensions")
	}
	_, _, err = image.Decode(bytes.NewReader(data))
	return err
}

func (s *authService) UpdateProfile(ctx context.Context, userID int, req *domain.UpdateProfileRequest) (*domain.User, error) {
	profile, err := validateProfile(*req)
	if err != nil {
		return nil, err
	}
	return s.userRepo.UpdateProfile(ctx, userID, &profile)
}
