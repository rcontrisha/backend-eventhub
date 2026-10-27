package controller

import (
	"fmt"
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"
	"rcontrisha/backend-eventhub/pkg"

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

// Get All Communities
//
// @Summary		Get All Communities
// @Description	Retrieve a paginated list of communities with optional search and filters
// @Tags		communities
// @Accept		json
// @Produce		json
// @Param		search	query	string	false	"Search communities by name"	example(Go community) maxLength(100)
// @Param		location	query	string	false	"Filter communities by location"	example(Jakarta) maxLength(100)
// @Param		tag		query	string	false	"Filter communities by tag"	example(technology) maxLength(50)
// @Param		page	query	int		false	"Page number (default: 1)"	minimum(1) example(1)
// @Param		limit	query	int		false	"Number of communities per page (default: 10, max: 100)"	minimum(1) maximum(100) example(10)
// @Router		/communities	[get]
// @Success		200	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		500	{object}	dto.Response
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

// Get Community Detail
//
// @Summary		Get Community Detail
// @Description	Retrieve details of a specific community
// @Tags		communities
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"Community ID (UUID)"	format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Router		/communities/{id}	[get]
// @Success		200	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		404	{object}	dto.Response
// @Failure		500	{object}	dto.Response
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

// Get Community Upcoming Events
//
// @Summary		Get Community Upcoming Events
// @Description	Retrieve upcoming events associated with a community
// @Tags		communities
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"Community ID (UUID)"	format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Router		/communities/{id}/upcoming-events	[get]
// @Success		200	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		500	{object}	dto.Response
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

// Get Community Members
//
// @Summary		Get Community Members
// @Description	Retrieve members of a community
// @Tags		communities
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"Community ID (UUID)"	format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Router		/communities/{id}/members	[get]
// @Success		200	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *CommunityController) GetCommunityMembers(ctx *gin.Context) {
	var req dto.GetEventDetailRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid community ID parameter",
			Data:    nil,
		})
		return
	}

	result, err := c.service.GetCommunityMembers(ctx, req.Id)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "error",
			Message: "Failed to fetch community members: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "community members retrieved successfully",
		Data: gin.H{
			"members": result,
		},
	})
}

// Get Popular Communities
//
// @Summary		Get Popular Communities
// @Description	Retrieve popular communities
// @Tags		communities
// @Accept		json
// @Produce		json
// @Router		/communities/popular	[get]
// @Success		200	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *CommunityController) GetPopularCommunities(ctx *gin.Context) {
	result, err := c.service.GetPopularCommunities(ctx)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "error",
			Message: "Failed to fetch communities: " + err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Communities retrieved successfully",
		Data: gin.H{
			"communities": result,
		},
	})
}

// Toggle Community Membership
//
// @Summary		Toggle Community Membership
// @Description	Join a community if not a member, or leave it if already joined
// @Tags		communities
// @Accept		json
// @Produce		json
// @Param		communityId	path	string	true	"Community ID (UUID)"	format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Router		/communities/{communityId}/toggle-join	[post]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *CommunityController) JoinOrLeaveController(ctx *gin.Context) {
	communityId := ctx.Param("communityId")

	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	userId := claims.Id
	action, err := c.service.JoinOrLeaveCommunity(ctx, communityId, userId)
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
		Message: fmt.Sprintf("Successfully %s the community.", action),
		Data: gin.H{
			"action":       action,
			"community_id": communityId,
		},
	})
}

// Get My Communities
//
// @Summary		Get My Communities
// @Description	Retrieve communities the authenticated user has joined
// @Tags		communities
// @Accept		json
// @Produce		json
// @Router		/communities/my-communities	[get]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		500	{object}	dto.Response
func (c *CommunityController) GetMyCommunities(ctx *gin.Context) {
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

	result, err := c.service.GetMyCommunities(ctx, userId)
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
		Message: "Events retrieved successfully",
		Data: gin.H{
			"communities": result,
		},
	})
}
