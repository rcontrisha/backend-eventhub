package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"rida@example.com"`
	Password string `json:"password" binding:"required" example:"StrongPass123!"`
}

type LoginResponse struct {
	Id        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string  `json:"email" example:"rida@example.com"`
	Name      string  `json:"name" example:"Rida"`
	AvatarUrl *string `json:"avatar_url" example:"https://example.com/avatar.png"`
	Role      string  `json:"role" example:"user"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"Rida"`
	Email    string `json:"email" binding:"required,email" example:"rida@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"StrongPass123!"`
}
