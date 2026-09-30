package controller

import (
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		service: *service,
	}
}

// Login
//
// @Summary				Login User
// @Description		Authenticate user and return token
// @Tags					auth
// @Accept				json
// @Produce				json
// @Router				/auth/login	[post]
// @Param					data	body	dto.LoginRequest	true	"login credentials"
// @Success				200		{object}	dto.Response
// @Failure				401		{object}	dto.Response
// @Failure				500		{object}	dto.Response
func (a *AuthController) LoginController(ctx *gin.Context) {
	var payload dto.LoginRequest
	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("error", e.Error())
		ctx.JSON(500, gin.H{
			"success": false,
			"msg":     e.Error(),
		})
		return
	}

	user, token, err := a.service.LoginService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Login Failed. Invalid email or password.",
			Data:    gin.H{},
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Login Success.",
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func (a *AuthController) RegisterController(ctx *gin.Context) {
	var payload dto.RegisterRequest
	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("Error: ", e.Error())
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: e.Error(),
			Data:    gin.H{},
		})
		return
	}

	if err := a.service.RegisterService(ctx.Request.Context(), payload); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}
	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Success Create New User.",
		Data:    gin.H{},
	})
}
