package router

import (
	"log"
	"rcontrisha/backend-eventhub/internal/controller"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/internal/service"
	"rcontrisha/backend-eventhub/pkg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRouter(r *gin.Engine, db *pgxpool.Pool) {
	authRouter := r.Group("/auth")

	repo := repository.NewAuthRepo(db)
	service := service.NewAuthService(repo)
	controller := controller.NewAuthController(service)

	authRouter.POST("/login", controller.LoginController)
	authRouter.POST("/register", controller.RegisterController)

	authRouter.POST("pwd", func(ctx *gin.Context) {
		type body struct {
			Password string `json:"pwd"`
		}
		var reqBody body
		if err := ctx.ShouldBindWith(&reqBody, binding.JSON); err != nil {
			log.Println("error", err.Error())
			// binding error
			ctx.JSON(500, gin.H{
				"success": false,
				"data":    nil,
				"msg":     "terjadi kesalahan server",
			})
			return
		}

		hc := pkg.NewRecommendedHashConfig()
		hash := hc.GenHash(reqBody.Password)

		ctx.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"pwd":  reqBody.Password,
				"hash": hash,
			},
		})
	})
}
