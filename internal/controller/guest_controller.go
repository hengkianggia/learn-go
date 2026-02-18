package controller

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/pagination"
	"learn/internal/pkg/request"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GuestController interface {
	CreateGuest(c *gin.Context)
	GetAllGuests(c *gin.Context)
	GetGuestBySlug(c *gin.Context)
	UpdateGuest(c *gin.Context)
}

type guestController struct {
	guestService service.GuestService
	logger       *slog.Logger
	db           *gorm.DB
}

func NewGuestController(guestService service.GuestService, logger *slog.Logger, db *gorm.DB) GuestController {
	return &guestController{guestService: guestService, logger: logger, db: db}
}

func (ctrl *guestController) CreateGuest(c *gin.Context) {
	var input dto.CreateGuestInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "create guest") {
		return
	}

	guest, err := ctrl.guestService.CreateGuest(input)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Guest created successfully", dto.ToGuestResponse(*guest))
}

func (ctrl *guestController) GetAllGuests(c *gin.Context) {
	var guests []model.Guest
	paginatedResponse, err := pagination.Paginate(c, ctrl.db, &model.Guest{}, &guests)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	if len(guests) == 0 {
		paginatedResponse.Data = make([]dto.GuestResponse, 0)
	} else {
		paginatedResponse.Data = dto.ToGuestResponses(guests)
	}

	response.SendSuccess(c, http.StatusOK, "Guests retrieved successfully", paginatedResponse)
}

func (ctrl *guestController) GetGuestBySlug(c *gin.Context) {
	slug := c.Param("slug")
	guest, err := ctrl.guestService.GetGuestBySlug(slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendNotFoundError(c, "Guest not found")
			return
		}
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, "Guest retrieved successfully", dto.ToGuestResponse(*guest))
}

func (ctrl *guestController) UpdateGuest(c *gin.Context) {
	slug := c.Param("slug")
	var input dto.UpdateGuestInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "update guest") {
		return
	}

	guest, err := ctrl.guestService.UpdateGuest(slug, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendNotFoundError(c, "Guest not found")
			return
		}
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, "Guest updated successfully", dto.ToGuestResponse(*guest))
}
