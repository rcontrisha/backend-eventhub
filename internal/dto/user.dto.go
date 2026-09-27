package dto

type UserProfile struct {
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url"`
	Location  *string `json:"location"`
	Bio       *string `json:"bio"`
	Role      string  `json:"role"`
}