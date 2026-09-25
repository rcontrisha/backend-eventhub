package repository

import (
	"context"
	"errors"
	"log"

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
	log.Printf("Payload - Repo: %s", email)
	sql := "SELECT id, email, password, name, avatar_url, location, bio, role FROM users WHERE email=$1"
	log.Printf("Query: %s", sql)

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
			log.Printf("LoginRepo: user not found for email %s", email)
			return model.User{}, err
		}

		log.Printf("LoginRepo database error: %v", err)
		return model.User{}, err
	}

	log.Printf("Users - Repo: %+v", user)
	return user, nil
}
