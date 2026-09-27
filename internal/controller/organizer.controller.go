package controller

import (
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"
	"rcontrisha/backend-eventhub/pkg"

	"github.com/gin-gonic/gin"
)

type OrganizerController struct {
	service *service.OrganizerService
}

func NewOrganizerController(service *service.OrganizerService) *OrganizerController {
	return &OrganizerController{
		service: service,
	}
}

func (o *OrganizerController) GetDashboard(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	userId := claims.Id
	data, err := o.service.GetDashboard(ctx, userId)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Successfully retrieve organizer's data.",
		Data: gin.H{
			"data":   data,
		},
	})
} 