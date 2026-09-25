package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Id        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url"`
	Location  *string `json:"location"`
	Bio       *string `json:"bio"`
	Role      string  `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTClaims(id, email, name string, avatar_url, location, bio *string, role string) *JWTClaims {
	return &JWTClaims{
		Id:        id,
		Email:     email,
		Name:      name,
		AvatarUrl: avatar_url,
		Location:  location,
		Bio:       bio,
		Role:      role,
		Issuer:    os.Getenv("JWT_ISSUER"),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
	}
}

func (j *JWTClaims) GenToken() (string, error) {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		return "", errors.New("JWT Key Not Found.")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func (j *JWTClaims) DecodeToken(token string) error {
	jwtToken, err := jwt.ParseWithClaims(token, j, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return err
	}
	if !jwtToken.Valid {
		return jwt.ErrTokenExpired
	}
	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}
	return nil
}
