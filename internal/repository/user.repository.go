package repository

import (
	"context"
	"errors"
	"fmt"
	"rcontrisha/backend-eventhub/internal/dto"
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
			u.role,
			u.created_at,
			u.updated_at
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
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &res, nil
}

func (u *UserRepo) ChangeUserProfile(ctx context.Context, userId string, payload dto.UserProfile) (*model.User, error) {
	query := `
		UPDATE users
		SET 
			name = COALESCE($1, name),
			avatar_url = COALESCE($2, avatar_url),
			location = COALESCE($3, location),
			bio = COALESCE($4, bio),
			updated_at = NOW()
		WHERE id = $5
		RETURNING id, email, name, avatar_url, location, bio, role, created_at, updated_at
	`
	
	args := []any{payload.Name, payload.AvatarUrl, payload.Location, payload.Bio, userId}
	var user model.User
	err := u.db.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&user.AvatarUrl,
		&user.Location,
		&user.Bio,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}
	
	return &user, nil
}