package controller

import (
	"net/http"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	service *service.EventService
}

func NewEventController(service *service.EventService) *EventController {
	return &EventController{
		service: service,
	}
}

func (e *EventController) GetAllEvents(ctx *gin.Context) {
	var req dto.GetEventsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  "error",
			Message: "Invalid query parameters: " + err.Error(),
			Data:    nil,
		})
		return
	}

	result, err := e.service.GetAllEvents(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  "error",
			Message: "Failed to fetch events: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  "success",
		Message: "Events retrieved successfully",
		Data: gin.H{
			"events": result.Events,
			"total":  result.Total,
			"page":   result.Page,
			"limit":  result.Limit,
		},
	})
}

func (e *EventController) GetEventDetail(ctx *gin.Context) {
	var req dto.GetEventDetailRequest

	// Bind URI param (contoh: /events/:id)
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  "error",
			Message: "Invalid event ID parameter",
			Data:    nil,
		})
		return
	}

	result, err := e.service.GetEventDetail(ctx.Request.Context(), req.Id)
	if err != nil {
		if err.Error() == "event not found" {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Status:  "error",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  "error",
			Message: "Failed to retrieve event details: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  "success",
		Message: "Event details retrieved successfully",
		Data: gin.H{
			"event": result,
		},
	})
}
