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

func (a *AuthRepo) FindAccount(ctx context.Context, email string) (model.User, error) {
	sql := "SELECT id, email, password, name, avatar_url, role FROM users WHERE email=$1"

	var user model.User

	err := a.db.QueryRow(ctx, sql, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.AvatarUrl,
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

func (a *AuthRepo) CreateNewUser(ctx context.Context, payload model.User) error {
	sql := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())"
	args := []any{payload.Name, payload.Email, payload.Password}
	cmd, err := a.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("No Row(s) Affected")
	}

	return nil
}
