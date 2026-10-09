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
			u.id,
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
		&res.Id,
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

func (u *UserRepo) GetUserPwd(ctx context.Context, userId string) (string, error) {
	query := `SELECT password FROM users WHERE id=$1`
	var oldPwd string
	if err := u.db.QueryRow(ctx, query, userId).Scan(&oldPwd); err != nil {
		return "", fmt.Errorf("failed to update user profile: %w", err)
	}

	return oldPwd, nil
}

func (u *UserRepo) ChangePassword(ctx context.Context, userId, newPwd string) error {
	query := `UPDATE users SET password=$1 WHERE id=$2`
	cmd, err := u.db.Exec(ctx, query, newPwd, userId)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("No Row(s) Affected")
	}

	return nil
}

func (u *UserRepo) GetUserInfo(ctx context.Context, userId string) (*model.UserInfo, error) {
	query := `
		SELECT
			u.id,
			u.email,
			u.name,
			u.avatar_url,
			u.bio,
			u.location,
			u.role,
			COALESCE((
				SELECT json_agg(ep.event_id)
				FROM public.event_participants ep
				WHERE ep.user_id = u.id
			), '[]'::json) AS joined_events,
			COALESCE((
				SELECT json_agg(se.event_id)
				FROM public.saved_events se
				WHERE se.user_id = u.id
			), '[]'::json) AS saved_events,
			COALESCE((
				SELECT json_agg(cm.community_id)
				FROM public.community_members cm
				WHERE cm.user_id = u.id
			), '[]'::json) AS joined_communities,
			u.created_at,
			u.updated_at
		FROM public.users u
		WHERE u.id = $1
	`

	var res model.UserInfo

	if err := u.db.QueryRow(ctx, query, userId).Scan(
		&res.Id,
		&res.Email,
		&res.Name,
		&res.AvatarUrl,
		&res.Bio,
		&res.Location,
		&res.Role,
		&res.JoinedEvents,
		&res.SavedEvents,
		&res.JoinedCommunities,
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
