package request

import (
	"learn/internal/pkg/response"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// BindJSONOrError attempts to bind the request body to the given input struct.
// If binding fails, it logs the error, sends a BadRequest response, and returns false.
// If binding succeeds, it returns true.
func BindJSONOrError(c *gin.Context, input interface{}, logger *slog.Logger, action string) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		logger.Warn("failed to bind JSON for "+action, slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return false
	}
	return true
}
