package controller

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/filters"
	"learn/internal/pkg/pagination"
	"learn/internal/pkg/request"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type venueController struct {
	venueService service.VenueService
	logger       *slog.Logger
	db           *gorm.DB
}

type VenueController interface {
	CreateVenue(c *gin.Context)
	GetAllVenues(c *gin.Context)
	GetVenueBySlug(c *gin.Context)
	UpdateVenue(c *gin.Context)
}

func NewVenueController(venueService service.VenueService, logger *slog.Logger, db *gorm.DB) VenueController {
	return &venueController{venueService: venueService, logger: logger, db: db}
}

func (ctrl *venueController) CreateVenue(c *gin.Context) {
	var input dto.CreateVenueInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "create venue") {
		return
	}

	venue, err := ctrl.venueService.CreateVenue(input)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Venue created successfully", dto.ToVenueResponse(*venue))
}

func (ctrl *venueController) GetAllVenues(c *gin.Context) {
	var venues []model.Venue

	filterFuncs := []filters.FilterFunc{
		filters.WithSearch("name"),
	}

	db := filters.ApplyFilter(ctrl.db, c, filterFuncs...)

	paginatedResponse, err := pagination.Paginate(c, db, &model.Venue{}, &venues)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	if len(venues) == 0 {
		paginatedResponse.Data = make([]dto.VenueResponse, 0)
	} else {
		paginatedResponse.Data = dto.ToVenueResponses(venues)
	}

	response.SendSuccess(c, http.StatusOK, "Venues retrieved successfully", paginatedResponse)
}

func (ctrl *venueController) GetVenueBySlug(c *gin.Context) {
	slug := c.Param("slug")
	venue, err := ctrl.venueService.GetVenueBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendNotFoundError(c, "Venue not found")
			return
		}
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, "Venue retrieved successfully", dto.ToVenueResponse(*venue))
}

func (ctrl *venueController) UpdateVenue(c *gin.Context) {
	slug := c.Param("slug")
	var input dto.UpdateVenueInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "update venue") {
		return
	}

	venue, err := ctrl.venueService.UpdateVenue(slug, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendNotFoundError(c, "Venue not found")
			return
		}
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, "Venue updated successfully", dto.ToVenueResponse(*venue))
}
