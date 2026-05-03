package controller

import (
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/pagination"
	"learn/internal/pkg/request"
	"learn/internal/pkg/response"
	"learn/internal/service"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type adminController struct {
	adminService service.AdminService
	logger       *slog.Logger
	db           *gorm.DB
}

type AdminController interface {
	ApproveUser(c *gin.Context)
	RejectUser(c *gin.Context)
	BlockUser(c *gin.Context)
	UnblockUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	ListUsers(c *gin.Context)
}

func NewAdminController(adminService service.AdminService, logger *slog.Logger, db *gorm.DB) AdminController {
	return &adminController{adminService: adminService, logger: logger, db: db}
}

func (ctrl *adminController) ApproveUser(c *gin.Context) {
	var input dto.AdminUserActionInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "approve user") {
		return
	}

	user, err := ctrl.adminService.ApproveUser(input.UserID)
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	ctrl.logger.Info("admin approved user", slog.Uint64("user_id", uint64(input.UserID)))
	response.SendSuccess(c, http.StatusOK, "Organizer approved successfully", dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserType:    user.UserType,
		IsVerified:  user.IsVerified,
		IsApproved:  user.IsApproved,
		IsBlocked:   user.IsBlocked,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

func (ctrl *adminController) RejectUser(c *gin.Context) {
	var input dto.AdminUserActionInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "reject user") {
		return
	}

	err := ctrl.adminService.RejectUser(input.UserID)
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	ctrl.logger.Info("admin rejected user", slog.Uint64("user_id", uint64(input.UserID)))
	response.SendSuccess(c, http.StatusOK, "Organizer rejected and deleted successfully", nil)
}

func (ctrl *adminController) BlockUser(c *gin.Context) {
	var input dto.AdminUserActionInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "block user") {
		return
	}

	user, err := ctrl.adminService.BlockUser(input.UserID)
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	ctrl.logger.Info("admin blocked user", slog.Uint64("user_id", uint64(input.UserID)))
	response.SendSuccess(c, http.StatusOK, "User blocked successfully", dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserType:    user.UserType,
		IsVerified:  user.IsVerified,
		IsApproved:  user.IsApproved,
		IsBlocked:   user.IsBlocked,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

func (ctrl *adminController) UnblockUser(c *gin.Context) {
	var input dto.AdminUserActionInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "unblock user") {
		return
	}

	user, err := ctrl.adminService.UnblockUser(input.UserID)
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	ctrl.logger.Info("admin unblocked user", slog.Uint64("user_id", uint64(input.UserID)))
	response.SendSuccess(c, http.StatusOK, "User unblocked successfully", dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserType:    user.UserType,
		IsVerified:  user.IsVerified,
		IsApproved:  user.IsApproved,
		IsBlocked:   user.IsBlocked,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

func (ctrl *adminController) DeleteUser(c *gin.Context) {
	var input dto.AdminUserActionInput
	if !request.BindJSONOrError(c, &input, ctrl.logger, "delete user") {
		return
	}

	err := ctrl.adminService.DeleteUser(input.UserID)
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	ctrl.logger.Info("admin deleted user", slog.Uint64("user_id", uint64(input.UserID)))
	response.SendSuccess(c, http.StatusOK, "User deleted successfully", nil)
}

func (ctrl *adminController) ListUsers(c *gin.Context) {
	var users []model.User

	db := ctrl.db
	if userType := c.Query("user_type"); userType != "" {
		db = db.Where("user_type = ?", userType)
	}

	paginatedResponse, err := pagination.Paginate(c, db, &model.User{}, &users)
	if err != nil {
		response.SendInternalServerError(c, ctrl.logger, err)
		return
	}

	if len(users) == 0 {
		paginatedResponse.Data = make([]dto.UserResponse, 0)
	} else {
		userResponses := make([]dto.UserResponse, len(users))
		for i, user := range users {
			userResponses[i] = dto.UserResponse{
				ID:          user.ID,
				Name:        user.Name,
				Email:       user.Email,
				PhoneNumber: user.PhoneNumber,
				UserType:    user.UserType,
				IsVerified:  user.IsVerified,
				IsApproved:  user.IsApproved,
				IsBlocked:   user.IsBlocked,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
			}
		}
		paginatedResponse.Data = userResponses
	}

	response.SendSuccess(c, http.StatusOK, "Users retrieved successfully", paginatedResponse)
}

func (ctrl *adminController) DeleteUserByParam(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.SendBadRequestError(c, "invalid user ID")
		return
	}

	err = ctrl.adminService.DeleteUser(uint(userID))
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	ctrl.logger.Info("admin deleted user", slog.Uint64("user_id", userID))
	response.SendSuccess(c, http.StatusOK, "User deleted successfully", nil)
}
