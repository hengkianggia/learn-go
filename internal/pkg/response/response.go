package response

import (
	"log/slog"

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

// --- Helper Functions ---

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

// --- Specific Error Helpers ---

func SendInternalServerError(c *gin.Context, logger *slog.Logger, err error) {
	logger.Error("internal server error", slog.String("error", err.Error()))
	sendError(c, 500, "INTERNAL_SERVER_ERROR", "An unexpected error occurred. Please try again later.", nil)
}

func SendValidationError(c *gin.Context, details []ValidationErrorDetail) {
	sendError(c, 400, "VALIDATION_ERROR", "Invalid input provided.", details)
}

func SendBadRequestError(c *gin.Context, message string) {
	sendError(c, 400, "BAD_REQUEST", message, nil)
}

func SendUnauthorizedError(c *gin.Context, message string) {
	sendError(c, 401, "UNAUTHORIZED", message, nil)
}

func SendNotFoundError(c *gin.Context, message string) {
	sendError(c, 404, "NOT_FOUND", message, nil)
}
