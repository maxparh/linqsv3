package domain

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	Phone        string    `json:"phone"`
}

// CreateUserRequest используется для регистрации
type CreateUserRequest struct {
	Email    string `json:"email"`
	Phone    string    `json:"phone,omitempty"`
	Password string `json:"password"`
}

// LoginRequest используется для входа
type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
