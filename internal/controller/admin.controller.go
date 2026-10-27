package controller

import (
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService *service.AdminService
}

func NewAdminController(adminService *service.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// Get Admin Dashboard Stats
//
// @Summary		Get Admin Dashboard Stats
// @Description	Retrieve platform-wide user, event, community, and attendance statistics
// @Tags		admin
// @Accept		json
// @Produce		json
// @Router		/admin/dashboard	[get]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		403	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (ac *AdminController) GetDashboardStats(ctx *gin.Context) {
	stats, err := ac.adminService.GetDashboardStats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: "Failed to retrieve admin dashboard stats",
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Successfully retrieved dashboard stats",
		Data: gin.H{
			"data": stats,
		},
	})
}
