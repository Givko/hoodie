package security_test

import (
	"testing"

	"github.com/givko/hoodie/internal/domain"
	"github.com/givko/hoodie/internal/security"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_Generate_Verify_CustomClaims(t *testing.T) {
	user := domain.User{
		Username: "test",
		IsAdmin:  true,
	}

	token, err := security.GenerateToken(&user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err.Error())
	}

	if token == "" {
		t.Fatalf("GenerateToken() token is empty")
	}

	verifyToken, err := security.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken() error = %v", err.Error())
	}

	if verifyToken == nil {
		t.Fatalf("VerifyToken() token is nil")
	}

	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("VerifyToken() claims is not jwt.MapClaims")
	}

	if claims["sub"] != user.Username {
		t.Fatalf("VerifyToken() sub = %v, want %v", claims["sub"], user.Username)
	}

	if claims["admin"] != user.IsAdmin {
		t.Fatalf("VerifyToken() admin = %v, want %v", claims["admin"], user.IsAdmin)
	}
}
