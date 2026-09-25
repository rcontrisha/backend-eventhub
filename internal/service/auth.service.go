package service

import (
	"context"
	"errors"
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/repository"
	"rcontrisha/backend-eventhub/pkg"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	repo *repository.AuthRepo
}

func NewAuthService(repo *repository.AuthRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (u *AuthService) LoginService(ctx context.Context, payload dto.LoginRequest) (dto.LoginResponse, string, error) {
	log.Printf("Payload - Service: %s", payload)
	result, err := u.repo.LoginRepo(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginResponse{}, "", errors.New("Incorrect Email or Password.	")
		}
		return dto.LoginResponse{}, "", err
	}
	log.Printf("Result - Service: %s", result)

	user := dto.LoginResponse{
		Id:        result.Id,
		Email:     result.Email,
		Name:      result.Name,
		AvatarUrl: result.AvatarUrl,
		Role:      result.Role,
	}

	log.Printf("Payload Password: %s\nDB Password: %s", payload.Password, result.Password)
	if err := pkg.Compare(payload.Password, result.Password); err != nil {
		return dto.LoginResponse{}, "", err
	}

	claims := pkg.NewJWTClaims(result.Id, result.Email, result.Name, result.AvatarUrl, result.Location, result.Bio, result.Role)
	token, err := claims.GenToken()

	return user, token, nil
}
