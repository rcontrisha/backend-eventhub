package router

import (
	_ "rcontrisha/backend-eventhub/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool) {
	router.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouter(router, db)
	EventRouter(router, db)
	CommunityRouter(router, db)
	UserRouter(router, db)
	OrganizerRouter(router, db)
}
