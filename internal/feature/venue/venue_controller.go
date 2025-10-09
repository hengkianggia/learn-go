package venue

import (
	"learn/internal/pkg/pagination"
	"learn/internal/pkg/response"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VenueController interface {
	CreateVenue(c *gin.Context)
	GetAllVenues(c *gin.Context)
}

type venueController struct {
	venueService VenueService
	logger       *slog.Logger
	db           *gorm.DB
}

func NewVenueController(venueService VenueService, logger *slog.Logger, db *gorm.DB) VenueController {
	return &venueController{venueService: venueService, logger: logger, db: db}
}

func (ctrl *venueController) CreateVenue(c *gin.Context) {
	var input CreateVenueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for create venue", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	venue, err := ctrl.venueService.CreateVenue(input)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Venue created successfully", venue)
}

func (ctrl *venueController) GetAllVenues(c *gin.Context) {
	var venues []Venue
	paginatedResponse, err := pagination.Paginate(c, ctrl.db, &Venue{}, &venues)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, "Venues retrieved successfully", paginatedResponse)
}