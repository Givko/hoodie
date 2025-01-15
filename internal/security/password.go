package security

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"

	"github.com/plamendelchev/hoodie/internal/domain"
	"golang.org/x/crypto/argon2"
)

var time = uint32(1)
var memory = uint32(2 * 1024)
var threads = uint8(1)

// HashPassword generates a new salt and hashes the password using Argon2id.
// It returns the encoded password and salt, or an error if the operation fails.
func HashPassword(password string) (*domain.Password, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return &domain.Password{}, err
	}

	keyLen := uint32(32)

	// Derive the key
	key := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)
	return domain.NewPassword(b64Key, b64Salt), nil
}

// ComparePassword verifies a password against the stored, encoded hash.
// It returns true if the password matches the hash, and false otherwise.
func ComparePassword(password string, b64HashedPwd string, b64Salt string) bool {
	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false
	}

	expectedKey, err := base64.RawStdEncoding.DecodeString(b64HashedPwd)
	if err != nil {
		return false
	}

	keyLen := uint32(len(expectedKey))

	// Derive a new key from the given password and compare
	newKey := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	if len(newKey) != len(expectedKey) {
		return false
	}

	return bytes.Equal(newKey, expectedKey)
}
