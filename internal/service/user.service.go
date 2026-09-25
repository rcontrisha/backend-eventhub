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

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) LoginService(ctx context.Context, payload dto.LoginRequest) (string, error) {
	log.Printf("Payload - Service: %s", payload)
	result, err := u.repo.LoginRepo(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("Incorrect Email or Password.	")
		}
		return "", err
	}
	log.Printf("Result - Service: %s", result)

	// response := dto.LoginResponse{
	// 	Id:        result.Id,
	// 	Email:     result.Email,
	// 	Name:      result.Name,
	// 	AvatarUrl: result.AvatarUrl,
	// 	Location:  result.Location,
	// 	Bio:       result.Bio,
	// 	Role:      result.Role,
	// }

	log.Printf("Payload Password: %s\nDB Password: %s", payload.Password, result.Password)
	if err := pkg.Compare(payload.Password, result.Password); err != nil {
		return "", err
	}

	claims := pkg.NewJWTClaims(result.Id, result.Email, result.Name, result.AvatarUrl, result.Location, result.Bio, result.Role)

	return claims.GenToken()
}
