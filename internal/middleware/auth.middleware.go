package middleware

import (
	"errors"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/pkg"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckToken(ctx *gin.Context) {
	bearer := ctx.GetHeader("Authorization")
	if bearer == "" {
		ctx.AbortWithStatusJSON(401, dto.Response{
			Status:  "failed",
			Message: "Please Log In First.",
		})
		return
	}
	result := strings.Split(bearer, " ")
	if len(result) != 2 {
		ctx.AbortWithStatusJSON(401, dto.Response{
			Status:  "failed",
			Message: "Invalid Bearer Token",
		})
		return
	}
	if result[0] != "Bearer" {
		ctx.AbortWithStatusJSON(401, dto.Response{
			Status:  "failed",
			Message: "Invalid Bearer Token",
		})
		return
	}

	var token pkg.JWTClaims
	err := token.DecodeToken(result[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			ctx.AbortWithStatusJSON(401, dto.Response{
				Status:  "failed",
				Message: "Invalid Token",
			})
			return
		}
		ctx.AbortWithStatusJSON(500, dto.Response{
			Status:  "failed",
			Message: "Internal Error.",
		})
		return
	}
	ctx.Set("token", token)
	ctx.Next()
}

func IsOrganizer(ctx *gin.Context) {
	token, exists := ctx.Get("token")
	if !exists {
		ctx.AbortWithStatusJSON(401, dto.Response{
			Status:  "failed",
			Message: "missing token",
		})
		return
	}
	// cek role nya
	t, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.AbortWithStatusJSON(401, dto.Response{
			Status:  "failed",
			Message: "invalid claims",
		})
		return
	}

	if t.Role != "organizer" {
		ctx.AbortWithStatusJSON(403, dto.Response{
			Status:  "failed",
			Message: "no privilege",
		})
		return
	}
	ctx.Next()
}
