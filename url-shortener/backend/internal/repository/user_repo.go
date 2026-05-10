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
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)	
}

type postgresUserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepo{db: db}
}

func (r *postgresUserRepo) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (email, password_hash, phone) 
        VALUES ($1, $2, $3) 
        RETURNING id, created_at
    `
    // Если phone пустой — передаём NULL (postgres поймёт)
    var phoneParam interface{} = user.Phone
    if user.Phone == "" {
        phoneParam = nil
    }
    
    err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, phoneParam).Scan(&user.ID, &user.CreatedAt)
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

func (r *postgresUserRepo) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
    var user domain.User
    query := `SELECT id, email, password_hash, created_at, phone FROM users WHERE phone = $1`
    err := r.db.QueryRowContext(ctx, query, phone).Scan(
        &user.ID, 
        &user.Email, 
        &user.PasswordHash, 
        &user.CreatedAt,
        &user.Phone, // ← добавляем сканирование phone
    )
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
