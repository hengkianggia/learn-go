package guest

import (
	"learn/internal/pkg/response"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GuestController interface {
	CreateGuest(c *gin.Context)
}

type guestController struct {
	guestService GuestService
	logger       *slog.Logger
}

func NewGuestController(guestService GuestService, logger *slog.Logger) GuestController {
	return &guestController{guestService: guestService, logger: logger}
}

func (ctrl *guestController) CreateGuest(c *gin.Context) {
	var input CreateGuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for create guest", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	guest, err := ctrl.guestService.CreateGuest(input)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Guest created successfully", guest)
}