package access

import (
	"sync"
)

var users = sync.Map{}

// UserRepository is the repository for user data.
type UserRepository struct{}
