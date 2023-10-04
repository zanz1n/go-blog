package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/user"
)

type UserAuthPayload struct {
	UserId     uuid.UUID     `json:"sub" validate:"required"`
	Email      string        `json:"email" validate:"required"`
	ExpiryDate int64         `json:"exp" validate:"required"`
	IssuedAt   int64         `json:"iat" validate:"required"`
	Role       user.UserRole `json:"role" validate:"required"`
}

func (u *UserAuthPayload) Validate() error {
	if err := validate.Struct(u); err != nil {
		return errors.ErrInvalidAuthToken
	}

	if time.Now().Unix() > u.ExpiryDate {
		return errors.ErrInvalidAuthToken
	}

	return nil
}

type AuthRepository interface {
	EncodeUserToken(data *UserAuthPayload) (string, error)
	CreateUserToken(info *user.User) (string, error)
	DecodeUserToken(payload string) (*UserAuthPayload, error)
	AuthUser(info *user.User, passwd string) (string, error)
}
