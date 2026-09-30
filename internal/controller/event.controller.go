package controller

import (
	"fmt"
	"log"
	"net/http"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"
	"rcontrisha/backend-eventhub/pkg"

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

// Event Detail
//
// @Summary				Get Event Detail
// @Description		Retrieve details of a specific event
// @Tags					events
// @Accept				json
// @Produce				json
// @Router				/events/{id}	[get]
// @Param					id	path	string	true	"Event ID"
// @Success				200		{object}	dto.Response
// @Failure				400		{object}	dto.Response
// @Failure				404		{object}	dto.Response
// @Failure				500		{object}	dto.Response
func (e *EventController) GetEventDetail(ctx *gin.Context) {
	var req dto.GetEventDetailRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid event ID parameter",
			Data:    nil,
		})
		return
	}

	result, err := e.service.GetEventDetail(ctx.Request.Context(), req.Id)
	if err != nil {
		if err.Error() == "event not found" {
			ctx.JSON(404, dto.Response{
				Status:  "error",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}

		ctx.JSON(500, dto.Response{
			Status:  "error",
			Message: "Failed to retrieve event details: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Event details retrieved successfully",
		Data: gin.H{
			"event": result,
		},
	})
}

func (e *EventController) JoinOrLeaveController(ctx *gin.Context) {
	eventId := ctx.Param("eventId")

	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	userId := claims.Id
	action, err := e.service.JoinOrLeaveEvent(ctx, eventId, userId)
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
		Message: fmt.Sprintf("Successfully %s the event.", action),
		Data: gin.H{
			"action":   action,
			"event_id": eventId,
		},
	})
}

func (e *EventController) GetUpcomingEvents(ctx *gin.Context) {
	result, err := e.service.GetUpcomingEvents(ctx)
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
			"events": result,
		},
	})
}

func (e *EventController) GetMyEvents(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	userId := claims.Id
	log.Println(userId)

	result, err := e.service.GetMyEvents(ctx, userId)
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
			"events": result,
		},
	})
}
