package router

import (
	"rcontrisha/backend-eventhub/internal/controller"
	"rcontrisha/backend-eventhub/internal/middleware"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AdminRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	adminRouter := r.Group("admin")

	repo := repository.NewAdminRepo(db)
	service := service.NewAdminService(repo)
	controller := controller.NewAdminController(service)

	adminRouter.GET("dashboard", middleware.CheckToken(rdb), middleware.IsAdmin, controller.GetDashboardStats)
}
