package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/sqli"
	"github.com/zanz1n/go-htmx/internal/user"
	"github.com/zanz1n/go-htmx/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func NewPostgresRepository(dba sqli.Querier) UserRepository {
	return &UserPostgresRepository{dba: dba}
}

type UserPostgresRepository struct {
	dba sqli.Querier
}

func (r *UserPostgresRepository) GetById(id uuid.UUID) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	u, err := r.dba.GetUserById(ctx, pguuid(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		} else {
			return nil, errors.ErrUserFetchFailed
		}
	}

	return pgToApiUser(u), nil
}

func (r *UserPostgresRepository) GetByEmail(email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	u, err := r.dba.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrUserNotFound
		} else {
			return nil, errors.ErrUserFetchFailed
		}
	}

	return pgToApiUser(u), nil
}

func (r *UserPostgresRepository) Create(data *user.UserCreateData) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := sqli.CreateUserParams{
		ID:       pguuid(uuid.New()),
		Username: data.Username,
		Email:    data.Email,
	}

	hash, err := bcrypt.GenerateFromPassword(utils.S2B(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrPasswordTooLong
	}

	user.Password = utils.B2S(hash)

	u, err := r.dba.CreateUser(ctx, &user)
	if err != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	return pgToApiUser(u), nil
}

func pgToApiUser(u *sqli.User) *user.User {
	return &user.User{
		ID:        uuid.UUID(u.ID.Bytes),
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		Role:      user.UserRole(u.Role),
	}
}

func pguuid(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}
