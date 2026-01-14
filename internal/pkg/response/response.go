package response

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Structs ---

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   struct {
		Code    string      `json:"code"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// --- Generic Helper Functions ---

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func sendError(c *gin.Context, statusCode int, code, message string, details interface{}) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	response.Error.Code = code
	response.Error.Details = details
	c.AbortWithStatusJSON(statusCode, response)
}

// --- Specific Success Helpers ---

func SendLoginSuccess(c *gin.Context, token string) {
	c.SetCookie("jwt_token", token, int(24*time.Hour/time.Second), "/", "localhost", false, true)
	SendSuccess(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

func SendLogoutSuccess(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	SendSuccess(c, http.StatusOK, "Logout successful", nil)
}

// --- Specific Error Helpers ---

func SendInternalServerError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Error("internal server error", slog.String("error", err.Error()))
	sendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred. Please try again later.", nil)
}

func SendValidationError(c *gin.Context, details []ValidationErrorDetail) {
	sendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input provided.", details)
}

func SendBadRequestError(c *gin.Context, message string) {
	sendError(c, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func SendUnauthorizedError(c *gin.Context, message string) {
	sendError(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func SendForbiddenError(c *gin.Context, message string) {
	sendError(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func SendNotFoundError(c *gin.Context, message string) {
	sendError(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}
