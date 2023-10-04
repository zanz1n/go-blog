package handlers

import (
	"github.com/zanz1n/go-htmx/internal/auth"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/user"
	ur "github.com/zanz1n/go-htmx/internal/user/repository"
)

func NewAuthHandlers(as auth.AuthRepository, us ur.UserRepository) *AuthHandlers {
	return &AuthHandlers{
		as: as,
		us: us,
	}
}

type AuthHandlers struct {
	as auth.AuthRepository
	us ur.UserRepository
}

func (h *AuthHandlers) HandleLogin(data *LoginIdenPayload) (string, error) {
	if err := validate.Struct(data); err != nil {
		return "", errors.ErrInvalidLoginData
	}

	u, err := h.us.GetByEmail(data.Email)
	if err != nil {
		e, ok := err.(*errors.StatusError)
		if !ok {
			e = errors.ErrUserFetchFailed
		}
		if e == errors.ErrUserNotFound {
			e = errors.ErrLoginFailed
		}

		return "", e
	}

	return h.as.AuthUser(u, data.Password)
}

func (h *AuthHandlers) HandleSignup(data *user.UserCreateData) (*user.User, string, error) {
	if err := validate.Struct(data); err != nil {
		return nil, "", errors.ErrInvalidSignupData
	}

	u, err := h.us.Create(data)
	if err != nil {
		return nil, "", err
	}

	token, err := h.as.CreateUserToken(u)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}
