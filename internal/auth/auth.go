package auth

import (
	"learn/internal/config"
	"learn/internal/middleware"
	"learn/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// GenerateJWT membuat token JWT baru
func GenerateJWT(user models.User) (string, error) {
	jwtKey := []byte(config.AppConfig.JWTSecretKey)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		Username: user.Username,
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
