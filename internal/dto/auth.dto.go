package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url"`
	Location  *string `json:"location"`
	Bio       *string `json:"bio"`
	Role      string  `json:"role"`
}
