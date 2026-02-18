package controller

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/request"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authController struct {
	authService service.AuthService
	logger      *slog.Logger
}

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	VerifyOTP(c *gin.Context)
	Profile(c *gin.Context)
	Logout(c *gin.Context)
}

func NewAuthController(authService service.AuthService, logger *slog.Logger) AuthController {
	return &authController{authService: authService, logger: logger}
}

func (ctrl *authController) Register(c *gin.Context) {
	var input dto.RegisterInput

	if !request.BindJSONOrError(c, &input, ctrl.logger, "registration") {
		return
	}

	user, err := ctrl.authService.Register(input)
	if err != nil {
		if err.Error() == "passwords do not match" {
			response.SendBadRequestError(c, "Passwords do not match")
		} else {
			response.SendBadRequestError(c, "A user with that email already exists")
		}
		return
	}

	ctrl.logger.Info("user registered successfully", slog.String("name", user.Name))

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
	if !request.BindJSONOrError(c, &input, ctrl.logger, "login") {
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

func (ctrl *authController) VerifyOTP(c *gin.Context) {
	var input dto.VerifyOTPInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "verify OTP") {
		return
	}

	err := ctrl.authService.VerifyOTP(input.Email, input.OTP)
	if err != nil {
		ctrl.logger.Warn("OTP verification failed", slog.String("email", input.Email), slog.String("error", err.Error()))
		response.SendBadRequestError(c, err.Error())
		return
	}

	ctrl.logger.Info("user verified successfully", slog.String("email", input.Email))
	response.SendSuccess(c, http.StatusOK, "Account verified successfully", nil)
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
