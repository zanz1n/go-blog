package user

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleCommon    UserRole = "COMMON"
	UserRolePublisher UserRole = "PUBLISHER"
	UserRoleAdmin     UserRole = "ADMIN"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      UserRole  `json:"role"`
}

type UserCreateData struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
