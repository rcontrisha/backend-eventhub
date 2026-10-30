package router

import (
	"rcontrisha/backend-eventhub/internal/broker"
	"rcontrisha/backend-eventhub/internal/controller"
	"rcontrisha/backend-eventhub/internal/middleware"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func EventRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client, brk *broker.Broker) {
	eventRouter := r.Group("/events")

	repo := repository.NewEventRepo(db)
	service := service.NewEventService(repo, brk)
	controller := controller.NewEventController(service, brk)

	eventRouter.GET("", controller.GetAllEvents)
	eventRouter.GET(":id", controller.GetEventDetail)
	eventRouter.GET("upcoming", controller.GetUpcomingEvents)
	eventRouter.GET("stream", controller.StreamAttendees)

	eventRouter.POST(":eventId/toggle-join", middleware.CheckToken(rdb), controller.JoinOrLeaveController)
	eventRouter.GET("my-events", middleware.CheckToken(rdb), controller.GetMyEvents)
	eventRouter.GET("saved-events", middleware.CheckToken(rdb), controller.GetSavedEvents)
}
