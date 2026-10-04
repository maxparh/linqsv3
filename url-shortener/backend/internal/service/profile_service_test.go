package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
)

type profileRepoStub struct {
	repository.UserRepository
	created   *domain.User
	updatedID int
	updated   *domain.UpdateProfileRequest
}

func (r *profileRepoStub) GetByEmail(context.Context, string) (*domain.User, error) {
	return nil, repository.ErrUserNotFound
}

func (r *profileRepoStub) Create(_ context.Context, user *domain.User) error {
	copy := *user
	r.created = &copy
	return nil
}

func (r *profileRepoStub) UpdateProfile(_ context.Context, id int, req *domain.UpdateProfileRequest) (*domain.User, error) {
	r.updatedID = id
	r.updated = req
	return &domain.User{ID: id, FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Phone: req.Phone}, nil
}

func TestProfileValidation(t *testing.T) {
	valid := domain.UpdateProfileRequest{FirstName: " Jane ", LastName: " Doe ", Email: " jane@example.com ", Phone: "8 (999) 123-45-67"}
	profile, err := validateProfile(valid)
	if err != nil || profile.FirstName != "Jane" || profile.LastName != "Doe" || profile.Email != "jane@example.com" || profile.Phone != "+79991234567" {
		t.Fatalf("unexpected normalized profile: %+v, %v", profile, err)
	}
	for _, phone := range []string{"", "+44 20 7946 0958"} {
		req := valid
		req.Phone = phone
		if _, err := validateProfile(req); err != nil {
			t.Fatalf("valid phone rejected: %v", err)
		}
	}
	for _, test := range []struct {
		name   string
		change func(*domain.UpdateProfileRequest)
	}{
		{"blank name", func(r *domain.UpdateProfileRequest) { r.FirstName = "  " }},
		{"long name", func(r *domain.UpdateProfileRequest) { r.LastName = strings.Repeat("a", 101) }},
		{"control character", func(r *domain.UpdateProfileRequest) { r.FirstName = "Ja\nne" }},
		{"email", func(r *domain.UpdateProfileRequest) { r.Email = "invalid" }},
		{"display email", func(r *domain.UpdateProfileRequest) { r.Email = "Jane <jane@example.com>" }},
		{"phone letters", func(r *domain.UpdateProfileRequest) { r.Phone = "+79991234567abc" }},
		{"short phone", func(r *domain.UpdateProfileRequest) { r.Phone = "+123" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := valid
			test.change(&req)
			if _, err := validateProfile(req); !errors.Is(err, ErrInvalidProfile) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}

func TestRegistrationPersistsProfile(t *testing.T) {
	repo := &profileRepoStub{}
	svc := NewAuthService(repo, "test-secret", 3600)
	user, err := svc.Register(context.Background(), &domain.CreateUserRequest{
		FirstName: " Jane ", LastName: " Doe ", Email: "jane@example.com", Phone: "+79991234567", Password: "test-password",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repo.created.FirstName != "Jane" || repo.created.LastName != "Doe" || repo.created.Phone != "+79991234567" {
		t.Fatalf("profile not persisted: %+v", repo.created)
	}
	if repo.created.PasswordHash == "" || user.PasswordHash != "" {
		t.Fatal("password hashing or response sanitization failed")
	}
}

func TestUpdateUsesAuthenticatedIDAndValidatesBeforeWrite(t *testing.T) {
	repo := &profileRepoStub{}
	svc := NewAuthService(repo, "test-secret", 3600)
	req := &domain.UpdateProfileRequest{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"}
	if _, err := svc.UpdateProfile(context.Background(), 42, req); err != nil {
		t.Fatal(err)
	}
	if repo.updatedID != 42 || repo.updated.Phone != "" {
		t.Fatal("wrong profile updated")
	}
	repo.updated = nil
	req.Email = "invalid"
	if _, err := svc.UpdateProfile(context.Background(), 42, req); !errors.Is(err, ErrInvalidProfile) {
		t.Fatal("invalid profile accepted")
	}
	if repo.updated != nil {
		t.Fatal("invalid profile written")
	}
}
