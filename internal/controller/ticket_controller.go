package controller

import (
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/request"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketController interface {
	CheckInTicket(c *gin.Context)
}

type ticketController struct {
	ticketService service.TicketService
	logger        *slog.Logger
}

func NewTicketController(ticketService service.TicketService, logger *slog.Logger) TicketController {
	return &ticketController{ticketService: ticketService, logger: logger}
}

func (ctrl *ticketController) CheckInTicket(c *gin.Context) {
	var input dto.CheckInTicketRequest
	if !request.BindJSONOrError(c, &input, ctrl.logger, "check in ticket") {
		return
	}

	userCtx, exists := c.Get("user")
	if !exists {
		response.SendUnauthorizedError(c, "User not authenticated")
		return
	}

	user, ok := userCtx.(model.User)
	if !ok {
		response.SendUnauthorizedError(c, "Invalid user context")
		return
	}

	result, err := ctrl.ticketService.CheckInTicket(input, user.ID)
	if err != nil {
		response.HandleAppError(c, err, ctrl.logger, "check in ticket")
		return
	}

	response.SendSuccess(c, http.StatusOK, "Ticket checked in successfully", result)
}
