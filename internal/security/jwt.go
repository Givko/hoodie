package security

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/plamendelchev/hoodie/internal/domain"
)

// GenerateToken generates a JWT token.
func GenerateToken(user *domain.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"admin":    user.IsAdmin,
	})

	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, nil
}
