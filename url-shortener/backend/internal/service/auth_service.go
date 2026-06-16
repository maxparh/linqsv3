package service

import (
    "context"
    "errors"
    "log"
    "strings"
    "time"
    "url-shortener/internal/domain"
    "url-shortener/internal/repository"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")


type AuthService interface {
	Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthTokens, error)
	ValidateToken(tokenString string) (*domain.Claims, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtSecret  []byte
	expiration time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, expirationSeconds int) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtSecret:  []byte(jwtSecret),
		expiration: time.Duration(expirationSeconds) * time.Second,
	}
}

func (s *authService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// Проверка, существует ли пользователь
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Хэширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
    Email:        req.Email,
    PasswordHash: string(hash),
    Phone:        req.Phone,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthTokens, error) {
	var user *domain.User
	var err error

	if strings.Contains(req.Identifier, "@") {
		// Поиск по email
		user, err = s.userRepo.GetByEmail(ctx, req.Identifier)
	} else {
		// Поиск по телефону
		normalizedPhone := normalizePhoneToE164(req.Identifier)
		user, err = s.userRepo.GetByPhone(ctx, normalizedPhone)
	}

	// Обработка ошибки
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	// Сравнение хэшей пароля (этот блок БЫЛ в оригинале, его нужно вернуть)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Генерация JWT
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Возврат токенов
	return &domain.AuthTokens{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.expiration.Seconds()),
	}, nil
}

func (s *authService) generateJWT(userID int, email string) (string, error) {
	claims := domain.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateToken(tokenString string) (*domain.Claims, error) {
	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("❌ ValidateToken: unexpected signing method: %v", token.Method)
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		log.Printf("❌ ValidateToken error: %v (token starts with: %s...)", err, tokenString[:min(30, len(tokenString))])
		return nil, err
	}

	if !token.Valid {
		log.Printf("❌ ValidateToken: token is not valid")
		return nil, errors.New("invalid token")
	}

	log.Printf("✅ ValidateToken: user_id=%d, email=%s", claims.UserID, claims.Email)
	return claims, nil
}

// Вспомогательная функция для Go <1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func normalizePhoneToE164(phone string) string {
	// Удаляем всё кроме цифр и +
	digits := ""
	for _, r := range phone {
		if r >= '0' && r <= '9' || r == '+' {
			digits += string(r)
		}
	}

	// Нормализация
	if strings.HasPrefix(digits, "8") && len(digits) == 11 {
		return "+7" + digits[1:]
	}
	if strings.HasPrefix(digits, "+7") {
		return digits
	}
	if strings.HasPrefix(digits, "7") && len(digits) == 11 {
		return "+" + digits
	}
	if len(digits) == 10 {
		return "+7" + digits
	}

	// Если не распознали — возвращаем как есть (поиск всё равно не найдёт)
	return digits
}