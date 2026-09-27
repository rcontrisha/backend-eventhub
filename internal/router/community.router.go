package router

import (
	"rcontrisha/backend-eventhub/internal/controller"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CommunityRouter(r *gin.Engine, db *pgxpool.Pool) {
	communityRouter := r.Group("/communities")

	repo := repository.NewCommunityRepo(db)
	service := service.NewCommunityService(repo)
	controller := controller.NewCommunityController(service)

	communityRouter.GET("", controller.GetAllCommunities)
}
