package event

import (
	"learn/internal/pkg/response"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventController interface {
	CreateEvent(c *gin.Context)
}

type eventController struct {
	eventService EventService
	logger       *slog.Logger
}

func NewEventController(eventService EventService, logger *slog.Logger) EventController {
	return &eventController{eventService: eventService, logger: logger}
}

func (ctrl *eventController) CreateEvent(c *gin.Context) {
	var input CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for create event", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	event, err := ctrl.eventService.CreateEvent(input)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Event created successfully", event)
}
