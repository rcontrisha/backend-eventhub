package router

import (
	"rcontrisha/backend-eventhub/internal/controller"
	"rcontrisha/backend-eventhub/internal/middleware"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRouter(r *gin.Engine, db *pgxpool.Pool) {
	userRouter := r.Group("/user")

	repo := repository.NewUserRepo(db)
	service := service.NewUserService(repo)
	controller := controller.NewUserController(service)

	userRouter.GET("profile", middleware.CheckToken, controller.GetUserProfile)
}
