package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
	Logout(c *gin.Context)
}

type authController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) AuthController {
	return &authController{authService: authService}
}

func (ctrl *authController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := ctrl.authService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (ctrl *authController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set JWT as an HTTP-only cookie
	// maxAge: 24 hours
	// path: /
	// domain: localhost (for development, change in production)
	// secure: false (for development, true in production for HTTPS)
	// httpOnly: true
	c.SetCookie("jwt_token", token, int(24*time.Hour/time.Second), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (ctrl *authController) Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "user": user})
}

func (ctrl *authController) Logout(c *gin.Context) {
	// Clear the jwt_token cookie by setting its MaxAge to -1 (expired)
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}