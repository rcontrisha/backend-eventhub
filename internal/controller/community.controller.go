package controller

import (
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	service *service.CommunityService
}

func NewCommunityController(service *service.CommunityService) *CommunityController {
	return &CommunityController{
		service: service,
	}
}

func (c *CommunityController) GetAllCommunities(ctx *gin.Context) {
	var req dto.GetCommunityRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid query parameters: " + err.Error(),
			Data:    nil,
		})
		return
	}

	result, err := c.service.GetAllCommunities(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "error",
			Message: "Failed to fetch events: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Communities retrieved successfully",
		Data: gin.H{
			"communities": result.Communities,
			"total":       result.Total,
			"page":        result.Page,
			"limit":       result.Limit,
		},
	})
}

func (c *CommunityController) GetCommunityDetail(ctx *gin.Context) {
	var req dto.GetEventDetailRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid community ID parameter",
			Data:    nil,
		})
		return
	}

	result, err := c.service.GetCommunityDetail(ctx, req.Id)
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
			Message: "Failed to retrieve community details: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Community details retrieved successfully",
		Data: gin.H{
			"community": result,
		},
	})
}

func (c *CommunityController) GetCommunityUpcomingEvents(ctx *gin.Context) {
	var req dto.GetEventDetailRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid community ID parameter",
			Data:    nil,
		})
		return
	}

	result, err := c.service.GetCommunityUpcomingEvents(ctx, req.Id)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "error",
			Message: "Failed to fetch community events: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "community events retrieved successfully",
		Data: gin.H{
			"events": result,
		},
	})
}
