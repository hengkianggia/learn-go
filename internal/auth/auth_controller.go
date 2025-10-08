package auth

import (
	"log/slog"
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
	logger      *slog.Logger
}

func NewAuthController(authService AuthService, logger *slog.Logger) AuthController {
	return &authController{authService: authService, logger: logger}
}

func (ctrl *authController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Error("failed to bind JSON for registration", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	_, err := ctrl.authService.Register(input)
	if err != nil {
		// Service layer already logs the specific error
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	ctrl.logger.Info("user registered successfully", slog.String("username", input.Username))
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (ctrl *authController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Error("failed to bind JSON for login", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token, err := ctrl.authService.Login(input)
	if err != nil {
		// Service layer already logs the specific error
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt_token", token, int(24*time.Hour/time.Second), "/", "localhost", false, true)

	ctrl.logger.Info("user logged in successfully", slog.String("username", input.Username))
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (ctrl *authController) Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		ctrl.logger.Error("user not found in context for profile access")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "user": user})
}

func (ctrl *authController) Logout(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	ctrl.logger.Info("user logged out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}