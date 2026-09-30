package controller

import (
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"
	"rcontrisha/backend-eventhub/pkg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
			"data": data,
		},
	})
}

func (o *OrganizerController) CreateEvent(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
	}
	claims := token.(pkg.JWTClaims)
	organizerId := claims.Id

	var payload dto.AddEventRequest
	if err := ctx.ShouldBindWith(&payload, binding.FormMultipart); err != nil {
		log.Println(err.Error())
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: "internal server error",
		})
		return
	}

	if err := o.service.AddEvent(ctx, organizerId, payload); err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Status:  "success",
		Message: "successfully create event.",
	})
}
