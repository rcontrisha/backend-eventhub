package router

import (
	"rcontrisha/backend-eventhub/internal/controller"
	"rcontrisha/backend-eventhub/internal/middleware"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func EventRouter(r *gin.Engine, db *pgxpool.Pool) {
	eventRouter := r.Group("/events")

	repo := repository.NewEventRepo(db)
	service := service.NewEventService(repo)
	controller := controller.NewEventController(service)

	eventRouter.GET("", controller.GetAllEvents)
	eventRouter.GET(":id", controller.GetEventDetail)

	eventRouter.POST(":eventId/toggle-join", middleware.CheckToken, controller.JoinOrLeaveController)
}
