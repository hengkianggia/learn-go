package auth

import (
	"errors"
	"learn/internal/pkg/response"
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
		ctrl.logger.Warn("failed to bind JSON for registration", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	user, err := ctrl.authService.Register(input)
	if err != nil {
		response.SendBadRequestError(c, "Username already exists")
		return
	}

	ctrl.logger.Info("user registered successfully", slog.String("username", user.Username))
	response.SendSuccess(c, http.StatusCreated, "Registration successful", gin.H{"id": user.ID, "username": user.Username})
}

func (ctrl *authController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for login", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	token, err := ctrl.authService.Login(input)
	if err != nil {
		response.SendUnauthorizedError(c, "Invalid username or password")
		return
	}

	c.SetCookie("jwt_token", token, int(24*time.Hour/time.Second), "/", "localhost", false, true)

	ctrl.logger.Info("user logged in successfully", slog.String("username", input.Username))
	response.SendSuccess(c, http.StatusOK, "Login successful", nil)
}

func (ctrl *authController) Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.SendInternalServerError(c, ctrl.logger, errors.New("user not found in context"))
		return
	}

	response.SendSuccess(c, http.StatusOK, "Profile retrieved successfully", user)
}

func (ctrl *authController) Logout(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	ctrl.logger.Info("user logged out successfully")
	response.SendSuccess(c, http.StatusOK, "Logout successful", nil)
}
