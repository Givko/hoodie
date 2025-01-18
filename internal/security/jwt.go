package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/plamendelchev/hoodie/internal/domain"
)

// GenerateToken generates a JWT token.
func GenerateToken(user *domain.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Username,
		"admin": user.IsAdmin,
		"iss":   "hoodie",
		"exp":   time.Now().UTC().Add(time.Second * 15).Unix(),
		"iat":   time.Now().UTC().Unix(),
	})

	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyToken verifies a JWT token.
func VerifyToken(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}
