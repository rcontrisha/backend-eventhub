package dto

import "time"

type UserProfile struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarUrl *string   `json:"avatar_url"`
	Location  *string   `json:"location"`
	Bio       *string   `json:"bio"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChangePassword struct {
	OldPwd string `json:"old_password"`
	NewPwd string `json:"new_password"`
}
