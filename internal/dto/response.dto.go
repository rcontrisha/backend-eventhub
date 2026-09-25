package dto

import "github.com/gin-gonic/gin"

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    gin.H `json:"data"`
}