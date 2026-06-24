package auth

import (
	"log/slog"

	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error to hash password", "error:", err)
		return "", exception.ErrSysHashPassword
	}

	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
