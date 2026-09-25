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
	Role      string  `json:"role"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
