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

// Get All Events
//
// @Summary		Get All Events
// @Description	Retrieve a paginated list of events with optional search and filters
// @Tags		events
// @Accept		json
// @Produce		json
// @Param		search	query	string	false	"Search events by title"	example(Go meetup) maxLength(100)
// @Param		location	query	string	false	"Filter events by location"	example(Jakarta) maxLength(100)
// @Param		tag		query	string	false	"Filter events by tag"	example(technology) maxLength(50)
// @Param		page	query	int		false	"Page number (default: 1)"	minimum(1) example(1)
// @Param		limit	query	int		false	"Number of events per page (default: 6, max: 100)"	minimum(1) maximum(100) example(6)
// @Router		/events	[get]
// @Success		200	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *EventController) GetAllEvents(ctx *gin.Context) {
	var req dto.GetEventsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  "error",
			Message: "Invalid query parameters: " + err.Error(),
			Data:    nil,
		})
		return
	}

	result, err := c.service.GetAllEvents(ctx.Request.Context(), req)
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

// Get Event Detail
//
// @Summary		Get Event Detail
// @Description	Retrieve details of a specific event
// @Tags		events
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"Event ID (UUID)"	format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Router		/events/{id}	[get]
// @Success		200	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		404	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *EventController) GetEventDetail(ctx *gin.Context) {
	var req dto.GetEventDetailRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid event ID parameter",
			Data:    nil,
		})
		return
	}

	result, err := c.service.GetEventDetail(ctx.Request.Context(), req.Id)
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

// Toggle Event Participation
//
// @Summary		Toggle Event Participation
// @Description	Join an event if not participating, or leave it if already joined
// @Tags		events
// @Accept		json
// @Produce		json
// @Param		eventId	path	string	true	"Event ID (UUID)"	format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Router		/events/{eventId}/toggle-join	[post]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *EventController) JoinOrLeaveController(ctx *gin.Context) {
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
	action, err := c.service.JoinOrLeaveEvent(ctx, eventId, userId)
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

// Get Upcoming Events
//
// @Summary		Get Upcoming Events
// @Description	Retrieve upcoming events
// @Tags		events
// @Accept		json
// @Produce		json
// @Router		/events/upcoming	[get]
// @Success		200	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *EventController) GetUpcomingEvents(ctx *gin.Context) {
	result, err := c.service.GetUpcomingEvents(ctx)
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

// Get My Events
//
// @Summary		Get My Events
// @Description	Retrieve events the authenticated user is participating in
// @Tags		events
// @Accept		json
// @Produce		json
// @Router		/events/my-events	[get]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *EventController) GetMyEvents(ctx *gin.Context) {
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

	result, err := c.service.GetMyEvents(ctx, userId)
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

// Get Saved Events
//
// @Summary		Get Saved Events
// @Description	Retrieve events saved by the authenticated user
// @Tags		events
// @Accept		json
// @Produce		json
// @Router		/events/saved-events	[get]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *EventController) GetSavedEvents(ctx *gin.Context) {
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

	result, err := c.service.GetSavedEvents(ctx, userId)
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
