package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/pkg"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	repo  *repository.UserRepo
	redis *redis.Client
}

func NewUserService(repo *repository.UserRepo, redis *redis.Client) *UserService {
	return &UserService{
		repo:  repo,
		redis: redis,
	}
}

func (u *UserService) GetUserProfile(ctx context.Context, userId string) (*dto.UserProfile, error) {
	key := fmt.Sprintf("rito:profile-%s", userId)

	if str, err := u.redis.Get(ctx, key).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			log.Println("key doesnt exist")
		} else {
			log.Println(err.Error())
		}
	} else {
		// cache hit
		var profile dto.UserProfile
		if err := json.Unmarshal([]byte(str), &profile); err != nil {
			log.Println("Parse error\nReason: ", err.Error())
		} else {
			return &profile, nil
		}
	}

	// cache miss
	result, err := u.repo.GetUserProfile(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := dto.UserProfile{
		Id:        result.Id,
		Email:     result.Email,
		Name:      result.Name,
		AvatarUrl: result.AvatarUrl,
		Location:  result.Location,
		Bio:       result.Bio,
		Role:      result.Role,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	str, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	if err := u.redis.Set(ctx, key, str, 0).Err(); err != nil {
		return nil, err
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

func (u *UserService) ChangeUserPwd(ctx context.Context, userId string, payload dto.ChangePassword) error {
	if userId == "" {
		return errors.New("unauthorized: missing user id")
	}

	oldPwd, err := u.repo.GetUserPwd(ctx, userId)
	if err := pkg.Compare(payload.OldPwd, oldPwd); err != nil {
		return errors.New("password mismatch")
	}

	if err != nil {
		return err
	}

	hashedNewPwd := pkg.NewRecommendedHashConfig().GenHash(payload.NewPwd)
	if err := u.repo.ChangePassword(ctx, userId, hashedNewPwd); err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserInfo(ctx context.Context, userId string) (*dto.UserInfo, error) {
	// key := fmt.Sprintf("rito:profile-%s", userId)

	// if str, err := u.redis.Get(ctx, key).Result(); err != nil {
	// 	if errors.Is(err, redis.Nil) {
	// 		log.Println("key doesnt exist")
	// 	} else {
	// 		log.Println(err.Error())
	// 	}
	// } else {
	// 	// cache hit
	// 	var profile dto.UserProfile
	// 	if err := json.Unmarshal([]byte(str), &profile); err != nil {
	// 		log.Println("Parse error\nReason: ", err.Error())
	// 	} else {
	// 		return &profile, nil
	// 	}
	// }

	// cache miss
	result, err := u.repo.GetUserInfo(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := dto.UserInfo{
		Id:                result.Id,
		Email:             result.Email,
		Name:              result.Name,
		AvatarUrl:         result.AvatarUrl,
		Location:          result.Location,
		Bio:               result.Bio,
		Role:              result.Role,
		JoinedEvents:      result.JoinedEvents,
		SavedEvents:       result.SavedEvents,
		JoinedCommunities: result.JoinedCommunities,
		CreatedAt:         result.CreatedAt,
		UpdatedAt:         result.UpdatedAt,
	}

	// str, err := json.Marshal(response)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := u.redis.Set(ctx, key, str, 0).Err(); err != nil {
	// 	return nil, err
	// }

	return &response, nil
}
