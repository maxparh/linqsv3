package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/domain"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

type profileServiceStub struct {
	service.AuthService
	id  int
	err error
}

func (s *profileServiceStub) GetProfile(_ context.Context, id int) (*domain.User, error) {
	s.id = id
	return &domain.User{ID: id, FirstName: "Jane", PasswordHash: "private-hash"}, s.err
}

func (s *profileServiceStub) UpdateProfile(_ context.Context, id int, _ *domain.UpdateProfileRequest) (*domain.User, error) {
	return s.GetProfile(context.Background(), id)
}

func TestProfileEndpointsUseContextIdentity(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPut} {
		t.Run(method, func(t *testing.T) {
			svc := &profileServiceStub{}
			h := NewAuthHandler(svc)
			handle := h.GetProfile
			if method == http.MethodPut {
				handle = h.UpdateProfile
			}
			r := httptest.NewRequest(method, "/api/profile?user_id=999", strings.NewReader(`{"first_name":"Jane","last_name":"Doe","email":"jane@example.com","phone":""}`))
			w := httptest.NewRecorder()
			handle(w, r)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("unauthenticated status: %d", w.Code)
			}
			r = r.WithContext(context.WithValue(r.Context(), middleware.UserIDKey, 42))
			w = httptest.NewRecorder()
			handle(w, r)
			if w.Code != http.StatusOK || svc.id != 42 {
				t.Fatalf("wrong identity/status: %d/%d", svc.id, w.Code)
			}
			if strings.Contains(w.Body.String(), "private-hash") || strings.Contains(w.Body.String(), "password") {
				t.Fatal("password exposed")
			}
		})
	}
}

func TestProfileRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	for _, body := range []string{`{"id":999}`, `{} {}`, `{`, `{"email":"` + strings.Repeat("x", 512*1024) + `"}`} {
		svc := &profileServiceStub{}
		r := httptest.NewRequest(http.MethodPut, "/api/profile", strings.NewReader(body))
		r = r.WithContext(context.WithValue(r.Context(), middleware.UserIDKey, 42))
		w := httptest.NewRecorder()
		NewAuthHandler(svc).UpdateProfile(w, r)
		if w.Code != http.StatusBadRequest || svc.id != 0 {
			t.Fatalf("invalid body accepted: %d", w.Code)
		}
	}
}

func TestProfileErrorStatuses(t *testing.T) {
	for _, test := range []struct {
		err    error
		status int
	}{
		{service.ErrInvalidProfile, http.StatusBadRequest},
		{repository.ErrEmailTaken, http.StatusConflict},
		{repository.ErrPhoneTaken, http.StatusConflict},
		{repository.ErrUserNotFound, http.StatusNotFound},
		{errors.New("private database details"), http.StatusInternalServerError},
	} {
		w := httptest.NewRecorder()
		writeProfileError(w, test.err)
		if w.Code != test.status {
			t.Fatalf("status: got %d want %d", w.Code, test.status)
		}
		if strings.Contains(w.Body.String(), "private database details") {
			t.Fatal("internal error exposed")
		}
	}
}
