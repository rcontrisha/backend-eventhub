package controller

import (
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: *service,
	}
}

func (u *UserController) LoginController(ctx *gin.Context) {
	var payload dto.LoginRequest
	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("error", e.Error())
		ctx.JSON(500, gin.H{
			"success": false,
			"msg":     e.Error(),
		})
		return
	}

	user, err := u.service.LoginService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, gin.H{
			"success": false,
			"message": "Login Failed. Invalid email or password.",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    user,
		"message": "Login Success.",
	})
}
