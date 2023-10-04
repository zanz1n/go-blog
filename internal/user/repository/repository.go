package repository

import (
	"github.com/google/uuid"
	"github.com/zanz1n/go-htmx/internal/user"
)

type UserRepository interface {
	GetById(id uuid.UUID) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Create(data *user.UserCreateData) (*user.User, error)
}
