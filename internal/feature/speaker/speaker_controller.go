package speaker

import (
	"learn/internal/pkg/response"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpeakerController interface {
	CreateSpeaker(c *gin.Context)
}

type speakerController struct {
	speakerService SpeakerService
	logger         *slog.Logger
}

func NewSpeakerController(speakerService SpeakerService, logger *slog.Logger) SpeakerController {
	return &speakerController{speakerService: speakerService, logger: logger}
}

func (ctrl *speakerController) CreateSpeaker(c *gin.Context) {
	var input CreateSpeakerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("failed to bind JSON for create speaker", slog.String("error", err.Error()))
		response.SendBadRequestError(c, "Invalid input format")
		return
	}

	speaker, err := ctrl.speakerService.CreateSpeaker(input)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, "Speaker created successfully", speaker)
}
