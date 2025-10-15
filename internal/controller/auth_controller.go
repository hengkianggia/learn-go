package controller

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
	Logout(c *gin.Context)
}

type authController struct {
	authService service.AuthService
	logger      *slog.Logger
}

func NewAuthController(authService service.AuthService, logger *slog.Logger) AuthController {
	return &authController{authService: authService, logger: logger}
}

func (ctrl *authController) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for registration", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	user, err := ctrl.authService.Register(input)
	if err != nil {
		if err.Error() == "passwords do not match" {
			response.SendBadRequestError(c, "Passwords do not match")
		} else {
			// This handles other errors like "email already exists"
			response.SendBadRequestError(c, "A user with that email already exists")
		}
		return
	}

	ctrl.logger.Info("user registered successfully", slog.String("name", user.Name))

	// Map user model to UserResponse DTO
	userResponse := dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserType:    user.UserType,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	response.SendSuccess(c, http.StatusCreated, "Registration successful", userResponse)
}

func (ctrl *authController) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for login", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	token, err := ctrl.authService.Login(input)
	if err != nil {
		response.SendUnauthorizedError(c, "Invalid email or password")
		return
	}

	ctrl.logger.Info("user logged in successfully", slog.String("email", input.Email))
	response.SendLoginSuccess(c, token)
}

func (ctrl *authController) Profile(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		response.SendInternalServerError(c, ctrl.logger, errors.New("user not found in context"))
		return
	}

	// Type assert user from context
	user, ok := userCtx.(model.User)
	if !ok {
		response.SendInternalServerError(c, ctrl.logger, errors.New("invalid user type in context"))
		return
	}

	// Map user model to UserResponse DTO
	userResponse := dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserType:    user.UserType,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	response.SendSuccess(c, http.StatusOK, "Profile retrieved successfully", userResponse)
}

func (ctrl *authController) Logout(c *gin.Context) {
	ctrl.logger.Info("user logged out successfully")
	response.SendLogoutSuccess(c)
}
