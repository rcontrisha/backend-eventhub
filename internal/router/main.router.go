package router

import (
	"path"
	_ "rcontrisha/backend-eventhub/docs"
	"rcontrisha/backend-eventhub/internal/broker"
	"rcontrisha/backend-eventhub/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client, brk *broker.Broker) {
	router.Use(middleware.Cors)

	router.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Static("img", path.Join("public", "img"))

	AuthRouter(router, db, rdb)
	EventRouter(router, db, rdb, brk)
	CommunityRouter(router, db, rdb)
	UserRouter(router, db, rdb)
	OrganizerRouter(router, db, rdb)
	AdminRouter(router, db, rdb)
}
