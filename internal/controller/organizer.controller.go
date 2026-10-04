package controller

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/service"
	"rcontrisha/backend-eventhub/pkg"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type OrganizerController struct {
	service *service.OrganizerService
}

func NewOrganizerController(service *service.OrganizerService) *OrganizerController {
	return &OrganizerController{
		service: service,
	}
}

func (o *OrganizerController) GetDashboard(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	userId := claims.Id
	data, err := o.service.GetDashboard(ctx, userId)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Status:  "success",
		Message: "Successfully retrieve organizer's data.",
		Data: gin.H{
			"data": data,
		},
	})
}

func (o *OrganizerController) CreateEvent(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
		return
	}
	claims := token.(pkg.JWTClaims)
	organizerId := claims.Id

	var payload dto.AddEventRequest
	if err := ctx.ShouldBindWith(&payload, binding.FormMultipart); err != nil {
		log.Println("[CreateEvent] Bind error:", err.Error())
		ctx.JSON(400, dto.Response{
			Status:  "failed",
			Message: err.Error(),
		})
		return
	}

	const maxFileSize = 4 * 1024 * 1024
	if payload.ImageUrl.Size > maxFileSize {
		ctx.JSON(400, dto.Response{
			Status:  "failed",
			Message: "file exceeds size limit (4mb)",
		})
		return
	}

	ext := strings.ToLower(path.Ext(payload.ImageUrl.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExtensions[ext] {
		ctx.JSON(400, dto.Response{
			Status:  "failed",
			Message: "format unsupported. please upload in .jpg, .jpeg, or .png format",
		})
		return
	}

	uploadDir := filepath.Join("public", "img", "events")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Println("[CreateEvent] Failed to create directory:", err.Error())
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: "error while processing file",
		})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), payload.ImageUrl.Filename)
	targetFilePath := filepath.Join(uploadDir, filename)

	if err := ctx.SaveUploadedFile(payload.ImageUrl, targetFilePath); err != nil {
		log.Println("[CreateEvent] SaveUploadedFile error:", err.Error())
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: "error saving file",
		})
		return
	}

	dbImageUrl := fmt.Sprintf("img/events/%s", filename)

	if err := o.service.AddEvent(ctx, organizerId, payload, dbImageUrl); err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "failed",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Status:  "success",
		Message: "successfully create event.",
	})
}

func (o *OrganizerController) EditEvent(ctx *gin.Context) {
	var req dto.EditEventRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid event ID parameter",
			Data:    nil,
		})
		return
	}

	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		ctx.JSON(400, dto.Response{
			Status:  "error",
			Message: "Invalid form data payload",
			Data:    nil,
		})
		return
	}

	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Status:  "failed",
			Message: "Token Data Not Found in Context.",
		})
		return
	}
	claims := token.(pkg.JWTClaims)
	organizerId := claims.Id

	var uploadedImageUrl *string
	if req.ImageUrl != nil {
		const maxFileSize = 2 * 1024 * 1024
		if req.ImageUrl.Size > maxFileSize {
			ctx.JSON(400, dto.Response{
				Status:  "failed",
				Message: "Ukuran gambar maksimal adalah 2MB",
			})
			return
		}

		ext := strings.ToLower(path.Ext(req.ImageUrl.Filename))
		allowedExtensions := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowedExtensions[ext] {
			ctx.JSON(400, dto.Response{
				Status:  "failed",
				Message: "unsupported file format. please upload .jpg, .jpeg, or .png format",
			})
			return
		}

		uploadDir := filepath.Join("public", "img", "events")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Println("[EditEvent] Failed to create directory:", err.Error())
			ctx.JSON(500, dto.Response{
				Status:  "failed",
				Message: "error processing file",
			})
			return
		}

		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), organizerId, ext)
		targetFilePath := filepath.Join(uploadDir, filename)

		if err := ctx.SaveUploadedFile(req.ImageUrl, targetFilePath); err != nil {
			log.Println("[EditEvent] SaveUploadedFile error:", err.Error())
			ctx.JSON(500, dto.Response{
				Status:  "failed",
				Message: "error saving file",
			})
			return
		}

		pathStr := fmt.Sprintf("img/events/%s", filename)
		uploadedImageUrl = &pathStr
	}

	if err := o.service.EditEvent(ctx, organizerId, req.Id, req, uploadedImageUrl); err != nil {
		ctx.JSON(500, dto.Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Event successfully updated",
	})
}
