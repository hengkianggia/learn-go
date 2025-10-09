package event

import (
	"learn/internal/pkg/pagination"
	"learn/internal/pkg/response"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventController interface {
	CreateEvent(c *gin.Context)
	GetAllEvents(c *gin.Context)
}

type eventController struct {
	eventService EventService
	logger       *slog.Logger
	db           *gorm.DB
}

func NewEventController(eventService EventService, logger *slog.Logger, db *gorm.DB) EventController {
	return &eventController{eventService: eventService, logger: logger, db: db}
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

	response.SendSuccess(c, http.StatusCreated, "Event created successfully", ToEventResponse(*event))
}

func (ctrl *eventController) GetAllEvents(c *gin.Context) {
	var events []Event
	db := ctrl.db.Preload("Venue").Preload("EventSpeakers.Speaker")
	paginatedResult, err := pagination.Paginate(c, db, &Event{}, &events)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	if len(events) == 0 {
		paginatedResult.Data = make([]EventResponse, 0)
	} else {
		paginatedResult.Data = ToEventResponses(events)
	}

	response.SendSuccess(c, http.StatusOK, "Events retrieved successfully", paginatedResult)
}
