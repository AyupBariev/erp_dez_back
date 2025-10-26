package services

import (
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	blacklist   *TokenBlacklist
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	blacklist *TokenBlacklist,
) *AuthService {
	expiry, _ := time.ParseDuration(os.Getenv("JWT_EXPIRE"))

	return &AuthService{
		userRepo:    userRepo,
		blacklist:   blacklist,
		jwtSecret:   []byte(os.Getenv("JWT_SECRET")),
		tokenExpiry: expiry,
	}
}

type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("find user error: %w", err)
	}
	if user == nil {
		return nil, nil // Пользователь не найден
	}

	// Сравнение хешированного пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil // Пароль не совпадает
	}

	return user, nil
}

func (s *AuthService) VerifyToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *AuthService) InvalidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return err
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}

	return s.blacklist.Add(tokenString, exp.Time)
}

func (s *AuthService) IsTokenValid(tokenString string) bool {
	// Проверяем в черном списке
	if invalid, _ := s.blacklist.Exists(tokenString); invalid {
		return false
	}

	// Проверяем подпись и срок действия
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	return err == nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}
	claims := JWTClaims{
		UserID: user.ID,
		Login:  user.Login,
		Role:   roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
