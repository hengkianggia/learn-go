package middleware

import (
	"errors"
	"learn/internal/config"
	"learn/internal/database"
	"learn/internal/model"
	"learn/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware adalah middleware untuk memproteksi route
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecretKey), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		var user model.User
		if err := database.DB.Where("email = ?", claims.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func RoleMiddleware(roles ...model.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, exists := c.Get("user")
		if !exists {
			response.SendUnauthorizedError(c, "User not found in context")
			c.Abort()
			return
		}

		user, ok := userCtx.(model.User)
		if !ok {
			response.SendInternalServerError(c, nil, errors.New("invalid user type in context"))
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if user.UserType == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.SendForbiddenError(c, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
