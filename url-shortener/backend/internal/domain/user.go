package domain

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Никогда не отдаем хеш в JSON
	CreatedAt    time.Time `json:"created_at"`
}

// CreateUserRequest используется для регистрации
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest используется для входа
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
