package controller

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/pagination"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventController interface {
	CreateEvent(c *gin.Context)
	GetAllEvents(c *gin.Context)
	GetEventBySlug(c *gin.Context)
	GetEventsByVenueSlug(c *gin.Context)
	GetEventsByGuestSlug(c *gin.Context)
}

type eventController struct {
	eventService service.EventService
	logger       *slog.Logger
	db           *gorm.DB
}

func NewEventController(eventService service.EventService, logger *slog.Logger, db *gorm.DB) EventController {
	return &eventController{eventService: eventService, logger: logger, db: db}
}

func (ctrl *eventController) CreateEvent(c *gin.Context) {
	var input dto.CreateEventInput
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

	response.SendSuccess(c, http.StatusCreated, "Event created successfully", dto.ToEventResponse(*event))
}

func (ctrl *eventController) GetAllEvents(c *gin.Context) {
	var events []model.Event
	db := ctrl.db.Preload("Venue").Preload("EventGuests.Guest").Preload("Prices")
	paginatedResult, err := pagination.Paginate(c, db, &model.Event{}, &events)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	if len(events) == 0 {
		paginatedResult.Data = make([]dto.EventResponse, 0)
	} else {
		paginatedResult.Data = dto.ToEventResponses(events)
	}

	response.SendSuccess(c, http.StatusOK, "Events retrieved successfully", paginatedResult)
}

func (ctrl *eventController) GetEventBySlug(c *gin.Context) {
	slug := c.Param("slug")
	event, err := ctrl.eventService.GetEventBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendNotFoundError(c, "Event not found")
			return
		}
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, "Event retrieved successfully", dto.ToEventResponse(*event))
}

func (ctrl *eventController) GetEventsByVenueSlug(c *gin.Context) {
	slug := c.Param("slug")
	var events []model.Event
	db := ctrl.db.Joins("JOIN venues ON venues.id = events.venue_id").Where("venues.slug = ?", slug).Preload("Venue").Preload("EventGuests.Guest").Preload("Prices")
	paginatedResult, err := pagination.Paginate(c, db, &model.Event{}, &events)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	if len(events) == 0 {
		paginatedResult.Data = make([]dto.EventResponse, 0)
	} else {
		paginatedResult.Data = dto.ToEventResponsesByVenue(events)
	}

	response.SendSuccess(c, http.StatusOK, "Events retrieved successfully", paginatedResult)
}

func (ctrl *eventController) GetEventsByGuestSlug(c *gin.Context) {
	slug := c.Param("slug")
	var events []model.Event
	db := ctrl.db.Joins("JOIN event_guests ON event_guests.event_id = events.id").Joins("JOIN guests ON guests.id = event_guests.guest_id").Where("guests.slug = ?", slug).Preload("Venue").Preload("EventGuests.Guest").Preload("Prices")
	paginatedResult, err := pagination.Paginate(c, db, &model.Event{}, &events)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	if len(events) == 0 {
		paginatedResult.Data = make([]dto.EventResponse, 0)
	} else {
		paginatedResult.Data = dto.ToEventResponses(events)
	}

	response.SendSuccess(c, http.StatusOK, "Events retrieved successfully", paginatedResult)
}
