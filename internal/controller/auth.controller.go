package controller

import (
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"
	"rcontrisha/backend-eventhub/pkg"

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

// Register
//
// @Summary				Register User
// @Description		Register a new user
// @Tags					auth
// @Accept				json
// @Produce				json
// @Router				/auth/register	[post]
// @Param					data	body	dto.RegisterRequest	true	"register data"
// @Success				200		{object}	dto.Response
// @Failure				500		{object}	dto.Response
func (a *AuthController) RegisterController(ctx *gin.Context) {
	var payload dto.RegisterRequest
	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("[Register Controller]-Binding Error: ", e.Error())
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: e.Error(),
			Data:    gin.H{},
		})
		return
	}

	if err := a.service.RegisterService(ctx.Request.Context(), payload); err != nil {
		log.Println("[Register Controller]-Service Error: ", err.Error())
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

// Logout
//
// @Summary				Logout User
// @Description		Logout the authenticated user and invalidate the token
// @Tags					auth
// @Accept				json
// @Produce				json
// @Router				/auth/logout	[post]
// @Security			BearerToken
// @Success				200		{object}	dto.Response
// @Failure				401		{object}	dto.Response
// @Failure				500		{object}	dto.Response
func (a *AuthController) Logout(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	jti := claims.ID
	expiresAt := claims.ExpiresAt
	if err := a.service.Logout(ctx, jti, expiresAt.Time); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Successfully logged out",
	})
}
