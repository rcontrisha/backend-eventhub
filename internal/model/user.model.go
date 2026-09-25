package model

import "time"

type User struct {
	Id        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	AvatarUrl *string    `db:"avatar_url"`
	Location  *string    `db:"location"`
	Bio       *string    `db:"bio"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// func NewUser(id, email, password, name, avatar_url, location, bio, role string, created_at, updated_at time.Time) *User {
// 	return &User{
// 		Id: id,
// 		Email: email,
// 		Password: password,
// 		Name: name,
// 		AvatarUrl: avatar_url,
// 		Location: location,
// 		Bio: bio,
// 		Role: role,
// 		CreatedAt: created_at,
// 		UpdatedAt: updated_at,
// 	}
// }
