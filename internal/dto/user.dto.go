package dto

import "time"

type UserProfile struct {
	Id        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string    `json:"email" example:"rida@example.com"`
	Name      string    `json:"name" example:"Rida"`
	AvatarUrl *string   `json:"avatar_url" example:"https://example.com/avatar.png"`
	Location  *string   `json:"location" example:"Jakarta, Indonesia"`
	Bio       *string   `json:"bio" example:"Community organizer and event enthusiast."`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2026-10-27T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-10-27T10:00:00Z"`
}

type UserInfo struct {
	Id        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string  `json:"email" example:"rida@example.com"`
	Name      string  `json:"name" example:"Rida"`
	AvatarUrl *string `json:"avatar_url" example:"https://example.com/avatar.png"`
	Location  *string `json:"location" example:"Jakarta, Indonesia"`
	Bio       *string `json:"bio" example:"Community organizer and event enthusiast."`
	Role      string  `json:"role" example:"user"`
	// JoinedEvents      []string  `json:"joined_events"`
	// SavedEvents       []string  `json:"saved_events"`
	// JoinedCommunities []string  `json:"joined_communities"`
	CreatedAt time.Time `json:"created_at" example:"2026-10-27T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-10-27T10:00:00Z"`
}

type ChangePassword struct {
	OldPwd string `json:"old_password" binding:"required" example:"CurrentPass123!"`
	NewPwd string `json:"new_password" binding:"required,min=8,nefield=OldPwd" example:"NewPass456!"`
}
