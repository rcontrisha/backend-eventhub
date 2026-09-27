package service

import (
	"context"
	"errors"
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
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	return &response, nil
}

func (u *UserService) ChangeUserProfile(ctx context.Context, userId string, payload dto.UserProfile) (*dto.UserProfile, error) {
	if userId == "" {
		return nil, errors.New("unauthorized: missing user id")
	}

	updatedUser, err := u.repo.ChangeUserProfile(ctx, userId, payload)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfile{
		Id:        updatedUser.Id,
		Email:     updatedUser.Email,
		Name:      updatedUser.Name,
		AvatarUrl: updatedUser.AvatarUrl,
		Role:      updatedUser.Role,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}