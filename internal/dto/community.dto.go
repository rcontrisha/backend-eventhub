package dto

import "time"

type GetCommunityRequest struct {
	Search   string `form:"search" binding:"omitempty,max=100" example:"Go community"`
	Location string `form:"location" binding:"omitempty,max=100" example:"Jakarta"`
	Tag      string `form:"tag" binding:"omitempty,max=50" example:"technology"`
	Page     int    `form:"page,default=1" binding:"omitempty,min=1" example:"1"`
	Limit    int    `form:"limit,default=10" binding:"omitempty,min=1,max=100" example:"10"`
}

type CommunityDetailResponse struct {
	Id             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           string    `json:"name" example:"Jakarta Go Community"`
	Description    string    `json:"description" example:"A community for Go developers in Jakarta."`
	BannerUrl      string    `json:"banner_url" example:"https://example.com/communities/go.png"`
	Tags           []string  `json:"tags" example:"technology,community"`
	MembersCount   int       `json:"members_count" example:"128"`
	UpcomingEvents int       `json:"upcoming_events" example:"4"`
	CreatedAt      time.Time `json:"created_at" example:"2026-10-27T10:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2026-10-27T10:00:00Z"`
}

type CommunityMembers struct {
	Name          string  `json:"name" example:"Rida"`
	AvatarUrl     *string `json:"avatar_url" example:"https://example.com/avatar.png"`
	CommunityRole string  `json:"role" example:"admin"`
}

type CommunityListResponse struct {
	Communities []CommunityListItemResponse `json:"communities"`
	Total       int                         `json:"total" example:"24"`
	Page        int                         `json:"page" example:"1"`
	Limit       int                         `json:"limit" example:"10"`
}

type CommunityListItemResponse struct {
	Id             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           string    `json:"name" example:"Jakarta Go Community"`
	Description    string    `json:"description" example:"A community for Go developers in Jakarta."`
	BannerUrl      string    `json:"banner_url" example:"https://example.com/communities/go.png"`
	Tags           []string  `json:"tags" example:"technology,community"`
	MembersCount   int       `json:"members_count" example:"128"`
	UpcomingEvents int       `json:"upcoming_events" example:"4"`
	CreatedAt      time.Time `json:"created_at" example:"2026-10-27T10:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2026-10-27T10:00:00Z"`
}
