package repository

import (
	"context"
	"errors"

	// "log"

	// "rcontrisha/backend-eventhub/internal/dto"
	"rcontrisha/backend-eventhub/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (u *AuthRepo) LoginRepo(ctx context.Context, email string) (model.User, error) {
	sql := "SELECT id, email, password, name, avatar_url, location, bio, role FROM users WHERE email=$1"

	var user model.User

	err := u.db.QueryRow(ctx, sql, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.AvatarUrl,
		&user.Location,
		&user.Bio,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, err
		}

		return model.User{}, err
	}

	return user, nil
}
