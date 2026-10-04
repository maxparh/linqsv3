package repository

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already in use")
var ErrPhoneTaken = errors.New("phone already in use")

type UserRepository interface {
	UpdateProfile(ctx context.Context, id int, req *domain.UpdateProfileRequest) (*domain.User, error)
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
        INSERT INTO users (email, password_hash, phone, first_name, last_name)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
	// Если phone пустой — передаём NULL (postgres поймёт)
	var phoneParam interface{} = user.Phone
	if user.Phone == "" {
		phoneParam = nil
	}

	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, phoneParam, user.FirstName, user.LastName).Scan(&user.ID, &user.CreatedAt)
	return userWriteError(err)
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
		&user.Phone,
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
	query := `SELECT id, email, created_at, COALESCE(phone, ''), first_name, last_name, avatar FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.Phone, &user.FirstName, &user.LastName, &user.Avatar)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepo) UpdateProfile(ctx context.Context, id int, req *domain.UpdateProfileRequest) (*domain.User, error) {
	var user domain.User
	query := `UPDATE users SET first_name = $1, last_name = $2, email = $3, phone = NULLIF($4, ''), avatar = COALESCE($6, avatar)
		WHERE id = $5
		RETURNING id, email, created_at, COALESCE(phone, ''), first_name, last_name, avatar`
	err := r.db.QueryRowContext(ctx, query, req.FirstName, req.LastName, req.Email, req.Phone, id, req.Avatar).
		Scan(&user.ID, &user.Email, &user.CreatedAt, &user.Phone, &user.FirstName, &user.LastName, &user.Avatar)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, userWriteError(err)
	}
	return &user, nil
}

func userWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return ErrEmailTaken
		case "users_phone_key":
			return ErrPhoneTaken
		}
	}
	return err
}
