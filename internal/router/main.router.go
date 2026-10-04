package router

import (
	_ "rcontrisha/backend-eventhub/docs"
	"rcontrisha/backend-eventhub/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.Use(middleware.Cors)

	router.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouter(router, db, rdb)
	EventRouter(router, db, rdb)
	CommunityRouter(router, db, rdb)
	UserRouter(router, db, rdb)
	OrganizerRouter(router, db, rdb)
	AdminRouter(router, db, rdb)
}
