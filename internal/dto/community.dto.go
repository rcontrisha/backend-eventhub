package dto

import "time"

type GetCommunityRequest struct {
	Search   string `form:"search" binding:"omitempty,max=100"`
	Location string `form:"location" binding:"omitempty,max=100"`
	Tag      string `form:"tag" binding:"omitempty,max=50"`
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit    int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

type CommunityDetailResponse struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	BannerUrl      string    `json:"banner_url"`
	Tags           []string  `json:"tags"`
	MembersCount   int       `json:"members_count"`
	UpcomingEvents int       `json:"upcoming_events"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CommunityMembers struct {
	Name          string  `json:"name"`
	AvatarUrl     *string `json:"avatar_url"`
	CommunityRole string  `json:"role"`
}

type CommunityListResponse struct {
	Communities []CommunityListItemResponse `json:"communities"`
	Total       int                         `json:"total"`
	Page        int                         `json:"page"`
	Limit       int                         `json:"limit"`
}

type CommunityListItemResponse struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	BannerUrl      string    `json:"banner_url"`
	Tags           []string  `json:"tags"`
	MembersCount   int       `json:"members_count"`
	UpcomingEvents int       `json:"upcoming_events"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
