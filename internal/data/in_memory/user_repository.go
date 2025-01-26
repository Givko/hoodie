package in_memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/givko/hoodie/internal/domain"
)

var users = sync.Map{}

type UserInMemoryRepository struct{}

// Add adds a new user to the repository.
// It returns an error if the user already exists.
func (repo UserInMemoryRepository) Add(user *domain.User) (err error) {
	user.CreatedAt = time.Now().UTC()
	user.IsAdmin = false
	_, isLoaded := users.LoadOrStore(user.Username, user)
	if isLoaded {
		return fmt.Errorf("user already exists")
	}

	return nil
}

// Get returns a user by username.
// It returns an error if the user does not exist or if user entry is corrupted.
func (repo UserInMemoryRepository) Get(username string) (dbValue *domain.User, err error) {
	value, ok := users.Load(username)
	if !ok {
		return &domain.User{}, fmt.Errorf("user not found")
	}

	user, ok := value.(*domain.User)
	if !ok {
		return &domain.User{}, fmt.Errorf("corrupted user entry")
	}

	return user, nil
}
