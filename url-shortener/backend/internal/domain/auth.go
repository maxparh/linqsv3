package domain

import "github.com/golang-jwt/jwt/v5"

type AuthTokens struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Claims struct {
	UserID               int    `json:"user_id"`
	Email                string `json:"email"`
	jwt.RegisteredClaims       
}