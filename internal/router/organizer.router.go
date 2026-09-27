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

	repo := repository.NewOrganizerRepo(db)
	service := service.NewOrganizerService(repo)
	controller := controller.NewOrganizerController(service)

	organizerRouter.GET("dashboard", middleware.CheckToken, middleware.IsOrganizer, controller.GetDashboard)
}
