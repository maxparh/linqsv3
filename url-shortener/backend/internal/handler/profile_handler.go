package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"url-shortener/internal/domain"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok || userID <= 0 {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}
	user, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		writeProfileError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok || userID <= 0 {
		WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}
	var req domain.UpdateProfileRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 512*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := h.authService.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		writeProfileError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, user)
}

func writeProfileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidProfile):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrEmailTaken), errors.Is(err, repository.ErrPhoneTaken):
		WriteError(w, http.StatusConflict, err.Error())
	case errors.Is(err, repository.ErrUserNotFound):
		WriteError(w, http.StatusNotFound, "user not found")
	default:
		log.Printf("Profile request failed: %v", err)
		WriteError(w, http.StatusInternalServerError, "could not process profile")
	}
}
