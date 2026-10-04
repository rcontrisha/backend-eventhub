package dto

import (
	"time"
)

type DashboardStatsResponse struct {
	TotalEvents    int     `json:"total_events"`
	TotalAttendees int     `json:"total_attendees"`
	AvgFillRate    float64 `json:"avg_fill_rate"`
	EventViews     int     `json:"event_views"`
}

type DashboardChartItem struct {
	Month string `json:"month"`
	Total int    `json:"total"`
}

type DashboardEventItemResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	ImageURL       *string   `json:"image_url"`
	Location       string    `json:"location"`
	StartTime      time.Time `json:"start_time"`
	Capacity       int       `json:"capacity"`
	AttendeesCount int       `json:"attendees_count"`
}

type DashboardUpcomingItemResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	StartTime      time.Time `json:"start_time"`
	Capacity       int       `json:"capacity"`
	AttendeesCount int       `json:"attendees_count"`
}

type OrganizerDashboardResponse struct {
	Stats          DashboardStatsResponse          `json:"stats"`
	ChartData      []DashboardChartItem            `json:"chart_data"`
	YourEvents     []DashboardEventItemResponse    `json:"your_events"`
	UpcomingEvents []DashboardUpcomingItemResponse `json:"upcoming_events"`
}

type AddEventDataRequest struct {
	Title       string    `form:"title" binding:"required,min=3,max=150"`
	Desc        string    `form:"desc" binding:"required,min=10"`
	ImageUrl    string    `form:"image_url" binding:"required"`
	CommunityId *string   `form:"community_id" binding:"omitempty,uuid"`
	StartTime   time.Time `form:"start_time" binding:"required"`
	EndTime     time.Time `form:"end_time" binding:"required,gtfield=StartTime"`
	Location    string    `form:"location" binding:"required,max=200"`
	Capacity    int       `form:"capacity" binding:"required,min=1"`
	Speakers    string    `form:"speakers" binding:"omitempty"`
}

type AddEventTagsRequest struct {
	Tags []string `form:"tags" binding:"omitempty,dive,uuid"`
}

type AddEventRequest struct {
	AddEventDataRequest
	AddEventTagsRequest
}

type EditEventDataRequest struct {
	Title       *string    `form:"title" binding:"omitempty,min=3,max=150"`
	Desc        *string    `form:"desc" binding:"omitempty,min=10"`
	ImageUrl    *string    `form:"image_url" binding:"omitempty"`
	CommunityId *string    `form:"community_id" binding:"omitempty,uuid"`
	StartTime   *time.Time `form:"start_time" binding:"omitempty"`
	EndTime     *time.Time `form:"end_time" binding:"omitempty,gtfield=StartTime"`
	Location    *string    `form:"location" binding:"omitempty,max=200"`
	Capacity    *int       `form:"capacity" binding:"omitempty,min=1"`
	Speakers    *string    `form:"speakers" binding:"omitempty"`
}

type EditEventTagsRequest struct {
	Tags []string `form:"tags" binding:"omitempty,dive,uuid"`
}

type EditEventRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
	EditEventDataRequest
	EditEventTagsRequest
}
