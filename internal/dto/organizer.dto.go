package dto

import (
	"mime/multipart"
	"time"
)

type DashboardStatsResponse struct {
	TotalEvents    int     `json:"total_events" example:"12"`
	TotalAttendees int     `json:"total_attendees" example:"480"`
	AvgFillRate    float64 `json:"avg_fill_rate" example:"72.5"`
	EventViews     int     `json:"event_views" example:"3200"`
}

type DashboardChartItem struct {
	Month string `json:"month" example:"2026-10"`
	Total int    `json:"total" example:"4"`
}

type DashboardEventItemResponse struct {
	ID             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title          string    `json:"title" example:"Go Community Meetup"`
	ImageURL       *string   `json:"image_url" example:"https://example.com/events/go-meetup.png"`
	Location       string    `json:"location" example:"Jakarta, Indonesia"`
	StartTime      time.Time `json:"start_time" example:"2026-11-15T09:00:00Z"`
	Capacity       int       `json:"capacity" example:"100"`
	AttendeesCount int       `json:"attendees_count" example:"42"`
}

type DashboardUpcomingItemResponse struct {
	ID             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title          string    `json:"title" example:"Go Community Meetup"`
	StartTime      time.Time `json:"start_time" example:"2026-11-15T09:00:00Z"`
	Capacity       int       `json:"capacity" example:"100"`
	AttendeesCount int       `json:"attendees_count" example:"42"`
}

type OrganizerDashboardResponse struct {
	Stats          DashboardStatsResponse          `json:"stats"`
	ChartData      []DashboardChartItem            `json:"chart_data"`
	YourEvents     []DashboardEventItemResponse    `json:"your_events"`
	UpcomingEvents []DashboardUpcomingItemResponse `json:"upcoming_events"`
}

type AddEventDataRequest struct {
	Title       string                `form:"title" binding:"required,min=3,max=150"`
	Desc        string                `form:"desc" binding:"required,min=10"`
	ImageUrl    *multipart.FileHeader `form:"image_url" binding:"required"`
	CommunityId *string               `form:"community_id" binding:"omitempty,uuid"`
	StartTime   time.Time             `form:"start_time" binding:"required"`
	EndTime     time.Time             `form:"end_time" binding:"required,gtfield=StartTime"`
	Location    string                `form:"location" binding:"required,max=200"`
	Capacity    int                   `form:"capacity" binding:"required,min=1"`
	Speakers    string                `form:"speakers" binding:"omitempty"`
}

type AddEventTagsRequest struct {
	Tags []string `form:"tags" binding:"omitempty,dive,uuid"`
}

type AddEventRequest struct {
	AddEventDataRequest
	AddEventTagsRequest
}

type EditEventDataRequest struct {
	Title       *string               `form:"title" binding:"omitempty,min=3,max=150"`
	Desc        *string               `form:"desc" binding:"omitempty,min=10"`
	ImageUrl    *multipart.FileHeader `form:"image_url" binding:"omitempty"`
	CommunityId *string               `form:"community_id" binding:"omitempty,uuid"`
	StartTime   *time.Time            `form:"start_time" binding:"omitempty"`
	EndTime     *time.Time            `form:"end_time" binding:"omitempty,gtfield=StartTime"`
	Location    *string               `form:"location" binding:"omitempty,max=200"`
	Capacity    *int                  `form:"capacity" binding:"omitempty,min=1"`
	Speakers    *string               `form:"speakers" binding:"omitempty"`
}

type EditEventTagsRequest struct {
	Tags []string `form:"tags" binding:"omitempty,dive,uuid"`
}

type EditEventRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
	EditEventDataRequest
	EditEventTagsRequest
}
