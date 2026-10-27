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

// Get Organizer Dashboard
//
// @Summary		Get Organizer Dashboard
// @Description	Retrieve dashboard statistics, event charts, and upcoming events for the authenticated organizer
// @Tags		organizer
// @Accept		json
// @Produce		json
// @Router		/organizer/dashboard	[get]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		403	{object}	dto.Response
// @Failure		500	{object}	dto.Response
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
		Data:    data,
	})
}

// Create Event
//
// @Summary		Create Event
// @Description	Create an event with an image upload. Image must be JPG, JPEG, or PNG and no larger than 4 MB.
// @Tags		organizer
// @Accept		multipart/form-data
// @Produce		json
// @Param		title			formData	string	true	"Event title"				minLength(3) maxLength(150) example(Go Community Meetup)
// @Param		desc			formData	string	true	"Event description"			minLength(10) example(A meetup for Go developers.)
// @Param		image_url		formData	file	true	"Event cover image (JPG, JPEG, or PNG; max 4 MB)"
// @Param		community_id	formData	string	false	"Community UUID"				format(uuid) example(770e8400-e29b-41d4-a716-446655440000)
// @Param		start_time		formData	string	true	"Event start time (RFC3339)"	format(date-time) example(2026-11-15T09:00:00Z)
// @Param		end_time		formData	string	true	"Event end time (RFC3339)"		format(date-time) example(2026-11-15T12:00:00Z)
// @Param		location		formData	string	true	"Event location"				maxLength(200) example(Jakarta, Indonesia)
// @Param		capacity		formData	int		true	"Maximum number of attendees"	minimum(1) example(100)
// @Param		speakers		formData	string	false	"Speakers data as a JSON string"	example([{"name":"Ayu Putri","role":"Software Engineer"}])
// @Param		tags			formData	[]string	false	"Tag UUIDs"
// @Router		/organizer/event	[post]
// @Security	BearerToken
// @Success		201	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		403	{object}	dto.Response
// @Failure		500	{object}	dto.Response
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

// Edit Event
//
// @Summary		Edit Event
// @Description	Update event fields using multipart form data. All fields are optional; uploaded images must be JPG, JPEG, or PNG and no larger than 4 MB.
// @Tags		organizer
// @Accept		multipart/form-data
// @Produce		json
// @Param		id				path		string	true	"Event UUID"						format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Param		title			formData	string	false	"Event title"						minLength(3) maxLength(150) example(Go Community Meetup)
// @Param		desc			formData	string	false	"Event description"					minLength(10) example(A meetup for Go developers.)
// @Param		image_url		formData	file	false	"Replacement cover image (JPG, JPEG, or PNG; max 4 MB)"
// @Param		community_id	formData	string	false	"Community UUID"						format(uuid) example(770e8400-e29b-41d4-a716-446655440000)
// @Param		start_time		formData	string	false	"Event start time (RFC3339)"			format(date-time) example(2026-11-15T09:00:00Z)
// @Param		end_time		formData	string	false	"Event end time (RFC3339)"				format(date-time) example(2026-11-15T12:00:00Z)
// @Param		location		formData	string	false	"Event location"						maxLength(200) example(Jakarta, Indonesia)
// @Param		capacity		formData	int		false	"Maximum number of attendees"			minimum(1) example(100)
// @Param		speakers		formData	string	false	"Speakers data as a JSON string"		example([{"name":"Ayu Putri","role":"Software Engineer"}])
// @Param		tags			formData	[]string	false	"Tag UUIDs"
// @Router		/organizer/event/{id}/edit	[patch]
// @Security	BearerToken
// @Success		200	{object}	dto.Response
// @Failure		400	{object}	dto.Response
// @Failure		401	{object}	dto.Response
// @Failure		403	{object}	dto.Response
// @Failure		500	{object}	dto.Response
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
		const maxFileSize = 4 * 1024 * 1024
		if req.ImageUrl.Size > maxFileSize {
			ctx.JSON(400, dto.Response{
				Status:  "failed",
				Message: "Ukuran gambar maksimal adalah 4MB",
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
