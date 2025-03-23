package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CompareHash compares a plaintext password with a hashed password
func CompareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword ensures password meets security requirements
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasNumber := strings.ContainsAny(password, "0123456789")
	return hasUpper && hasLower && hasNumber
}
