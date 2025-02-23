package in_memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/givko/hoodie/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserInMemoryRepository struct {
	redisClient *redis.Client
}

func NewUserInMemoryRepository(rd *redis.Client) *UserInMemoryRepository {
	return &UserInMemoryRepository{
		redisClient: rd,
	}
}

// Add adds a new user to the repository.
// It returns an error if the user already exists.
func (repo UserInMemoryRepository) Add(user *domain.User) (err error) {
	user.CreatedAt = time.Now().UTC()
	user.IsAdmin = false
	ctx := context.Background()

	errorMessage := repo.redisClient.JSONSetMode(ctx, user.Username, "$", user, "NX").Err()
	if errorMessage != nil {
		return errorMessage
	}

	return nil
}

// Get returns a user by username.
// It returns an error if the user does not exist or if user entry is corrupted.
func (repo UserInMemoryRepository) Get(username string) (dbValue *domain.User, err error) {
	ctx := context.Background()
	value, err := repo.redisClient.JSONGet(ctx, username).Result()
	if err != nil {
		return &domain.User{}, fmt.Errorf("user not found")
	}

	user := &domain.User{}
	err = json.Unmarshal([]byte(value), &user)
	if err != nil {
		return &domain.User{}, fmt.Errorf("corrupted user entry")
	}

	return user, nil
}
