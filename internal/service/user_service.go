package service

import (
	"github.com/google/uuid"
	"github.com/plamendelchev/hoodie/internal/api/dto"
	"github.com/plamendelchev/hoodie/internal/domain"
	"github.com/plamendelchev/hoodie/internal/security"
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
func (service *UserService) RegisterUser(createUser dto.CreateUser) error {
	hashedPassword, err := security.HashPassword(createUser.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Id:       uuid.New().String(),
		Username: createUser.Username,
		Password: *hashedPassword,
		IsAdmin:  false,
	}

	return service.userRepository.Add(&user)
}
