// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: queries.sql

package sqli

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPost = `-- name: CreatePost :one
INSERT INTO "posts" ("id", "title", "content", "topics", "description", "thumb_image", "user_id") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at, title, content, topics, description, thumb_image, user_id
`

type CreatePostParams struct {
	ID          pgtype.UUID
	Title       string
	Content     []byte
	Topics      []byte
	Description string
	ThumbImage  pgtype.UUID
	UserID      pgtype.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg *CreatePostParams) (*Post, error) {
	row := q.db.QueryRow(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.Topics,
		arg.Description,
		arg.ThumbImage,
		arg.UserID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Content,
		&i.Topics,
		&i.Description,
		&i.ThumbImage,
		&i.UserID,
	)
	return &i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO "users" ("id", "username", "email", "password") VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, email, username, password, role
`

type CreateUserParams struct {
	ID       pgtype.UUID
	Username string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) (*User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
	)
	return &i, err
}

const getPartialPostById = `-- name: GetPartialPostById :one
SELECT
    "posts"."id",
    "posts"."created_at",
    "posts"."updated_at",
    "posts"."title",
    "posts"."description",
    "posts"."thumb_image",
    "posts"."user_id",
    "users"."username" AS "user_username"
FROM "posts"
    INNER JOIN "users" ON "posts"."user_id" = "users"."id"
WHERE "posts"."id" = $1
`

type GetPartialPostByIdRow struct {
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
	Title        string
	Description  string
	ThumbImage   pgtype.UUID
	UserID       pgtype.UUID
	UserUsername string
}

func (q *Queries) GetPartialPostById(ctx context.Context, id pgtype.UUID) (*GetPartialPostByIdRow, error) {
	row := q.db.QueryRow(ctx, getPartialPostById, id)
	var i GetPartialPostByIdRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
		&i.ThumbImage,
		&i.UserID,
		&i.UserUsername,
	)
	return &i, err
}

const getPostById = `-- name: GetPostById :one
SELECT id, created_at, updated_at, title, content, topics, description, thumb_image, user_id FROM "posts" WHERE "id" = $1
`

func (q *Queries) GetPostById(ctx context.Context, id pgtype.UUID) (*Post, error) {
	row := q.db.QueryRow(ctx, getPostById, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Content,
		&i.Topics,
		&i.Description,
		&i.ThumbImage,
		&i.UserID,
	)
	return &i, err
}

const getPostByIdWithUser = `-- name: GetPostByIdWithUser :one
SELECT
    "posts"."id",
    "posts"."created_at",
    "posts"."updated_at",
    "posts"."title",
    "posts"."content",
    "posts"."topics",
    "posts"."description",
    "posts"."thumb_image",
    "posts"."user_id",
    "users"."created_at" AS "user_created_at",
    "users"."updated_at" AS "user_updated_at",
    "users"."email" AS "user_email",
    "users"."username" AS "user_username",
    "users"."role" AS "user_role"
FROM "posts"
    INNER JOIN "users" ON "posts"."user_id" = "users"."id"
WHERE "posts"."id" = $1
`

type GetPostByIdWithUserRow struct {
	ID            pgtype.UUID
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
	Title         string
	Content       []byte
	Topics        []byte
	Description   string
	ThumbImage    pgtype.UUID
	UserID        pgtype.UUID
	UserCreatedAt pgtype.Timestamp
	UserUpdatedAt pgtype.Timestamp
	UserEmail     string
	UserUsername  string
	UserRole      UserRole
}

func (q *Queries) GetPostByIdWithUser(ctx context.Context, id pgtype.UUID) (*GetPostByIdWithUserRow, error) {
	row := q.db.QueryRow(ctx, getPostByIdWithUser, id)
	var i GetPostByIdWithUserRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Content,
		&i.Topics,
		&i.Description,
		&i.ThumbImage,
		&i.UserID,
		&i.UserCreatedAt,
		&i.UserUpdatedAt,
		&i.UserEmail,
		&i.UserUsername,
		&i.UserRole,
	)
	return &i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, username, password, role FROM "users" WHERE "email" = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
	)
	return &i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, updated_at, email, username, password, role FROM "users" WHERE "id" = $1
`

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (*User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
	)
	return &i, err
}
