package domain

import "time"

type User struct {
	Avatar       string    `json:"avatar"`
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	Phone        string    `json:"phone"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
}

// CreateUserRequest используется для регистрации
type CreateUserRequest struct {
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateProfileRequest struct {
	Avatar    *string `json:"avatar,omitempty"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
}

// LoginRequest используется для входа
type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
