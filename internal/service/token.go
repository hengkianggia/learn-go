package service

import (
	"learn/internal/config"
	"learn/internal/middleware"
	"learn/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// GenerateJWT membuat token JWT baru
func GenerateJWT(user model.User) (string, error) {
	jwtKey := []byte(config.AppConfig.JWTSecretKey)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidatePassword membandingkan password yang di-hash dengan password plain text
func ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
