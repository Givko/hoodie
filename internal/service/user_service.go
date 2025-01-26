package service

import (
	"fmt"
	"time"

	"github.com/givko/hoodie/internal/api/contracts"
	"github.com/givko/hoodie/internal/domain"
	"github.com/givko/hoodie/internal/security"
	"github.com/google/uuid"
)

type UserRepository interface {
	Add(user *domain.User) error
	Get(username string) (*domain.User, error)
}

// UserService is a service that handles user-related operations.
type UserService struct {
	userRepository UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repository UserRepository) *UserService {
	return &UserService{userRepository: repository}
}

// RegisterUser registers a new user.
func (service *UserService) RegisterUser(createUser *contracts.CreateUser) error {
	hashedPassword, err := security.HashPassword(createUser.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Id:        uuid.New().String(),
		Username:  createUser.Username,
		Password:  *hashedPassword,
		IsAdmin:   false,
		CreatedAt: time.Now().UTC(),
	}

	return service.userRepository.Add(&user)
}

func (service *UserService) LoginUser(loginUser *contracts.LoginUser) (string, error) {
	user, err := service.userRepository.Get(loginUser.Username)
	if err != nil {
		return "", err
	}

	if !security.ComparePassword(loginUser.Password, user.Password.Password, user.Password.Salt) {
		return "", fmt.Errorf("invalid password")
	}

	jwt, err := security.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
