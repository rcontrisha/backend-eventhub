package dto

import "time"

type GetEventsRequest struct {
	Search   string `form:"search" binding:"omitempty,max=100" example:"Go meetup"`
	Location string `form:"location" binding:"omitempty,max=100" example:"Jakarta"`
	Tag      string `form:"tag" binding:"omitempty,max=50" example:"technology"`
	Page     int    `form:"page,default=1" binding:"omitempty,min=1" example:"1"`
	Limit    int    `form:"limit,default=6" binding:"omitempty,min=1,max=100" example:"6"`
}

type EventListItemResponse struct {
	Id             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title          string    `json:"title" example:"Go Community Meetup"`
	ImageURL       string    `json:"image_url" example:"https://example.com/events/go-meetup.png"`
	Tags           []string  `json:"tags" example:"technology,community"`
	Capacity       *int      `json:"capacity" example:"100"`
	AttendeesCount int       `json:"attendees_count" example:"42"`
	StartTime      time.Time `json:"start_time" example:"2026-11-15T09:00:00Z"`
	EndTime        time.Time `json:"end_time" example:"2026-11-15T12:00:00Z"`
	Location       string    `json:"location" example:"Jakarta, Indonesia"`
	CreatedAt      time.Time `json:"created_at" example:"2026-10-27T10:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2026-10-27T10:00:00Z"`
}

type EventListResponse struct {
	Events []EventListItemResponse `json:"events" example:"[]"`
	Total  int                     `json:"total" example:"24"`
	Page   int                     `json:"page" example:"1"`
	Limit  int                     `json:"limit" example:"6"`
}

type GetEventDetailRequest struct {
	Id string `uri:"id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type EventOrganizerResponse struct {
	ID        string  `json:"id" example:"660e8400-e29b-41d4-a716-446655440000"`
	Name      string  `json:"name" example:"Rida"`
	Email     string  `json:"email" example:"organizer@example.com"`
	AvatarURL *string `json:"avatar_url" example:"https://example.com/organizer.png"`
	Role      string  `json:"role" example:"organizer"`
}

type EventCommunityResponse struct {
	ID          string  `json:"id" example:"770e8400-e29b-41d4-a716-446655440000"`
	Name        string  `json:"name" example:"Jakarta Go Community"`
	Description string  `json:"description" example:"A community for Go developers in Jakarta."`
	BannerURL   *string `json:"banner_url" example:"https://example.com/communities/go.png"`
}

type Speaker struct {
	Name string `json:"name" example:"Ayu Putri"`
	Role string `json:"role" example:"Software Engineer"`
}

type TagResponse struct {
	ID   string `json:"id" example:"880e8400-e29b-41d4-a716-446655440000"`
	Name string `json:"name" example:"Technology"`
}

type EventDetailResponse struct {
	Id                string                  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title             string                  `json:"title" example:"Go Community Meetup"`
	Desc              string                  `json:"desc" example:"A meetup for Go developers to share ideas and learn together."`
	ImageUrl          string                  `json:"image_url" example:"https://example.com/events/go-meetup.png"`
	Tags              []TagResponse           `json:"tags"`
	Organizer         EventOrganizerResponse  `json:"organizer"`
	Community         *EventCommunityResponse `json:"community"`
	Location          string                  `json:"location" example:"Jakarta, Indonesia"`
	StartTime         time.Time               `json:"start_time" example:"2026-11-15T09:00:00Z"`
	EndTime           time.Time               `json:"end_time" example:"2026-11-15T12:00:00Z"`
	Capacity          int                     `json:"capacity" example:"100"`
	TotalParticipants int                     `json:"total_participants" example:"42"`
	Speakers          []Speaker               `json:"speakers"`
	CreatedAt         time.Time               `json:"created_at" example:"2026-10-27T10:00:00Z"`
	UpdatedAt         time.Time               `json:"updated_at" example:"2026-10-27T10:00:00Z"`
}

type JoinOrLeaveEventResponse struct {
	EventId string `json:"event_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status  string `json:"status" example:"joined"`
}
