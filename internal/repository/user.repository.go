package repository

import (
	"context"
	"errors"
	"rcontrisha/backend-eventhub/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) GetUserProfile(ctx context.Context, userId string) (*model.User, error) {
	query := `
		SELECT
			u.email,
			u.name,
			u.avatar_url,
			u.bio,
			u.location,
			u.role
		FROM public.users u
		WHERE u.id = $1
	`

	var res model.User

	if err := u.db.QueryRow(ctx, query, userId).Scan(
		&res.Email,
		&res.Name,
		&res.AvatarUrl,
		&res.Bio,
		&res.Location,
		&res.Role,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &res, nil
}