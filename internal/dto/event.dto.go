package dto

import "time"

type GetEventsRequest struct {
	Search   string `form:"search"`   
	Location string `form:"location"` 
	Tag      string `form:"tag"`      
	Page     int    `form:"page,default=1"`
	Limit    int    `form:"limit,default=10"`
}

type EventListItemResponse struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	ImageURL       string    `json:"image_url"`
	Tags           []string  `json:"tags"`
	Capacity       *int      `json:"capacity"`
	AttendeesCount int       `json:"attendees_count"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Location       string    `json:"location"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type EventListResponse struct {
	Events []EventListItemResponse `json:"events"`
	Total  int                     `json:"total"`
	Page   int                     `json:"page"`
	Limit  int                     `json:"limit"`
}

type GetEventDetailRequest struct {
	Id string `uri:"id"`
}

type EventOrganizerResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url"`
	Role      string  `json:"role"`
}

type EventCommunityResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BannerURL   *string `json:"banner_url"`
}

type SpeakerResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type TagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EventDetailResponse struct {
	Id                string                  `json:"id"`
	Title             string                  `json:"title"`
	Desc              string                  `json:"desc"`
	ImageUrl          string                  `json:"image_url"`
	Tags              []TagResponse           `json:"tags"`
	Organizer         EventOrganizerResponse  `json:"organizer"`
	Community         *EventCommunityResponse `json:"community"`
	Location          string                  `json:"location"`
	StartTime         time.Time               `json:"start_time"`
	EndTime           time.Time               `json:"end_time"`
	Capacity          int                     `json:"capacity"`
	TotalParticipants int                     `json:"total_participants"`
	Speakers          []SpeakerResponse       `json:"speakers"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}

type JoinOrLeaveEventResponse struct {
	EventId string `json:"event_id"`
	Status string `json:"status"`
}