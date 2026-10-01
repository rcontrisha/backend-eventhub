package router

import (
	_ "rcontrisha/backend-eventhub/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouter(router, db)
	EventRouter(router, db)
	CommunityRouter(router, db)
	UserRouter(router, db, rdb)
	OrganizerRouter(router, db)
}
