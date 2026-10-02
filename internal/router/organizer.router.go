package router

import (
	"rcontrisha/backend-eventhub/internal/controller"
	"rcontrisha/backend-eventhub/internal/middleware"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OrganizerRouter(r *gin.Engine, db *pgxpool.Pool) {
	organizerRouter := r.Group("organizer")

	repo := repository.NewOrganizerRepo()
	service := service.NewOrganizerService(repo, db)
	controller := controller.NewOrganizerController(service)

	organizerRouter.GET("dashboard", middleware.CheckToken, middleware.IsOrganizer, controller.GetDashboard)
	organizerRouter.POST("event", middleware.CheckToken, middleware.IsOrganizer, controller.CreateEvent)
	organizerRouter.PATCH("event/:id/edit", middleware.CheckToken, middleware.IsOrganizer, controller.EditEvent)
}
