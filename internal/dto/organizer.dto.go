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
	Title       string    `form:"title"`
	Desc        string    `form:"desc"`
	ImageUrl    string    `form:"image_url"`
	CommunityId *string    `form:"community_id"`
	StartTime   time.Time `form:"start_time"`
	EndTime     time.Time `form:"end_time"`
	Location    string    `form:"location"`
	Capacity    int       `form:"capacity"`
	Speakers    string    `form:"speakers"`
}

type AddEventTagsRequest struct {
	Tags []string `form:"tags"`
}

type AddEventRequest struct {
	AddEventDataRequest
	AddEventTagsRequest
}
