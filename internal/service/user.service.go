package service

import (
	"context"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) GetUserProfile(ctx context.Context, userId string) (*dto.UserProfile, error) {
	result, err := u.repo.GetUserProfile(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := dto.UserProfile{
		Email: result.Email,
		Name: result.Name,
		AvatarUrl: result.AvatarUrl,
		Location: result.Location,
		Bio: result.Bio,
		Role: result.Role,
	}

	return &response, nil
}