package services

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetUserByLogin(login string) (*models.User, error) {
	user, err := s.userRepository.FindByUsername(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by login=%s: %w", login, err)
	}
	return user, nil
}

func (s *UserService) CreateUser(login, password, firstName, secondName string, roleID int64) (*models.User, error) {
	existingUser, err := s.userRepository.FindByUsername(login)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &models.User{
		FirstName:  firstName,
		SecondName: secondName,
		Login:      login,
		Password:   string(hashedPassword),
		RoleID:     roleID,
	}

	if err := s.userRepository.Create(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}
