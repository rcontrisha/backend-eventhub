package middleware

import (
	"errors"
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/pkg"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func CheckToken(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		blacklistKey := "rito:blacklist:" + token.ID
		// log.Println("Key to check: ", blacklistKey)
		exists, err := rdb.Exists(ctx, blacklistKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Println(err)
			ctx.AbortWithStatusJSON(500, dto.Response{
				Status:  "failed",
				Message: "Internal server error while checking token validity",
				Data:    nil,
			})
			return
		}

		if exists > 0 {
			ctx.AbortWithStatusJSON(401, dto.Response{
				Status:  "failed",
				Message: "Token has been revoked or logged out",
				Data:    nil,
			})
			return
		}

		ctx.Set("token", token)
		ctx.Next()
	}
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
