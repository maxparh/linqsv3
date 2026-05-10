package repository

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

type postgresUserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *postgresUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
