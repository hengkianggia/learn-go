package response

import (
	apperrors "learn/internal/errors"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// HandleAppError handles different types of application errors and sends appropriate HTTP responses
func HandleAppError(c *gin.Context, err error, logger *slog.Logger, operation string) bool {
	switch appErr := err.(type) {
	case apperrors.ValidationError:
		logger.Info("Validation error in "+operation,
			slog.String("field", appErr.Field),
			slog.String("message", appErr.Message),
			slog.Any("value", appErr.Value))
		SendBadRequestError(c, appErr.Error())
		return true

	case apperrors.BusinessRuleError:
		logger.Info("Business rule error in "+operation,
			slog.String("rule", appErr.Rule),
			slog.String("message", appErr.Message))
		SendBadRequestError(c, appErr.Error())
		return true

	case apperrors.SystemError:
		logger.Error("System error in "+operation,
			slog.String("operation", appErr.Operation),
			slog.String("message", appErr.Message),
			slog.Any("error", appErr.Err))
		SendInternalServerError(c, logger, appErr)
		return true

	default:
		// For any other error types, treat as internal server error
		logger.Error("Unknown error in "+operation, slog.String("error", err.Error()))
		SendInternalServerError(c, logger, err)
		return true
	}
}

// HandleAppErrorWithNotFound handles different types of application errors and sends appropriate HTTP responses
// This version sends a 404 error for BusinessRuleError instead of 400
func HandleAppErrorWithNotFound(c *gin.Context, err error, logger *slog.Logger, operation string) bool {
	switch appErr := err.(type) {
	case apperrors.ValidationError:
		logger.Info("Validation error in "+operation,
			slog.String("field", appErr.Field),
			slog.String("message", appErr.Message),
			slog.Any("value", appErr.Value))
		SendBadRequestError(c, appErr.Error())
		return true

	case apperrors.BusinessRuleError:
		logger.Info("Business rule error in "+operation,
			slog.String("rule", appErr.Rule),
			slog.String("message", appErr.Message))
		SendNotFoundError(c, appErr.Error())
		return true

	case apperrors.SystemError:
		logger.Error("System error in "+operation,
			slog.String("operation", appErr.Operation),
			slog.String("message", appErr.Message),
			slog.Any("error", appErr.Err))
		SendInternalServerError(c, logger, appErr)
		return true

	default:
		// For any other error types, treat as internal server error
		logger.Error("Unknown error in "+operation, slog.String("error", err.Error()))
		SendInternalServerError(c, logger, err)
		return true
	}
}
