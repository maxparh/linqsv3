package service

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"
	"testing"
	"url-shortener/internal/domain"
)

func TestAvatarValidation(t *testing.T) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 16, 16))); err != nil {
		t.Fatal(err)
	}
	valid := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes())
	if err := validateAvatar(valid); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []string{
		"https://example.com/photo.png",
		"data:image/svg+xml;base64,PHN2Zy8+",
		"data:image/png;base64,bm90LWFuLWltYWdl",
		strings.Replace(valid, "image/png", "image/jpeg", 1),
		"data:image/png;base64," + strings.Repeat("A", 400000),
		"data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes()[:30]),
	} {
		if err := validateAvatar(invalid); err == nil {
			t.Fatal("invalid avatar accepted")
		}
	}
	buffer.Reset()
	if err := png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 1025, 1))); err != nil {
		t.Fatal(err)
	}
	if err := validateAvatar("data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes())); err == nil {
		t.Fatal("oversized dimensions accepted")
	}
}

func TestAvatarCanBeOmittedOrRemoved(t *testing.T) {
	req := domain.UpdateProfileRequest{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"}
	if _, err := validateProfile(req); err != nil {
		t.Fatal(err)
	}
	empty := ""
	req.Avatar = &empty
	if _, err := validateProfile(req); err != nil {
		t.Fatal(err)
	}
	invalid := "invalid"
	req.Avatar = &invalid
	if _, err := validateProfile(req); err == nil {
		t.Fatal("invalid avatar accepted")
	}
}
