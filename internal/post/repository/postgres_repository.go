package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/post"
	"github.com/zanz1n/go-htmx/internal/sqli"
	"github.com/zanz1n/go-htmx/internal/user"
	"github.com/zanz1n/go-htmx/internal/utils"
)

func NewPostgresRepository(dba sqli.Querier) PostRepository {
	return &PostPostgresRepository{
		dba: dba,
	}
}

type PostPostgresRepository struct {
	dba sqli.Querier
}

func (r *PostPostgresRepository) GetById(id uuid.UUID) (*post.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := r.dba.GetPostById(ctx, pguuid(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrPostNotFound
		}
		return nil, errors.ErrPostFetchFailed
	}

	return pgToApiPost(data), nil
}

func (r *PostPostgresRepository) GetByIdWithUser(id uuid.UUID) (*post.PostWithUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := r.dba.GetPostByIdWithUser(ctx, pguuid(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrPostNotFound
		}
		return nil, errors.ErrPostFetchFailed
	}

	return pgToApiPostWithUser(data), nil
}

func (r *PostPostgresRepository) GetPartialById(id uuid.UUID) (*post.PartialPost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := r.dba.GetPartialPostById(ctx, pguuid(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrPostNotFound
		}
		return nil, errors.ErrPostFetchFailed
	}

	return pgToApiPartialPost(data), nil
}

func (r *PostPostgresRepository) Create(userId uuid.UUID, data *post.PostCreateData) (*post.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if data.Description == nil {
		data.Description = str("")
	}

	dbdata, err := r.dba.CreatePost(ctx, &sqli.CreatePostParams{
		ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Title:       data.Title,
		Content:     utils.S2B(data.Content),
		Topics:      []byte("[]"),
		Description: *data.Description,
		ThumbImage:  pguuid(data.ThumbImage),
		UserID:      pguuid(userId),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrPostNotFound
		}
		return nil, errors.ErrPostFetchFailed
	}

	return pgToApiPost(dbdata), nil
}

func pgToApiPartialPost(p *sqli.GetPartialPostByIdRow) *post.PartialPost {
	return &post.PartialPost{
		ID:           p.ID.Bytes,
		UserId:       p.UserID.Bytes,
		UserUsername: p.UserUsername,
		ThumbImage:   apiuuid(p.ThumbImage),
		CreatedAt:    p.CreatedAt.Time,
		UpdatedAt:    p.UpdatedAt.Time,
		Title:        p.Title,
		Description:  p.Description,
	}
}

func pgToApiPostWithUser(p *sqli.GetPostByIdWithUserRow) *post.PostWithUser {
	topics := []post.PostTopic{}
	if err := json.Unmarshal(p.Topics, &topics); err != nil {
		slog.Error("Failed to unmarshal post topics", "error", err)
		topics = []post.PostTopic{}
	}

	return &post.PostWithUser{
		Post: post.Post{
			ID:          p.ID.Bytes,
			UserId:      p.UserID.Bytes,
			ThumbImage:  apiuuid(p.ThumbImage),
			CreatedAt:   p.CreatedAt.Time,
			UpdatedAt:   p.UpdatedAt.Time,
			Title:       p.Title,
			Description: p.Description,
			Content:     utils.B2S(p.Content),
			Topics:      topics,
		},
		User: user.User{
			ID:        p.UserID.Bytes,
			CreatedAt: p.UserCreatedAt.Time,
			UpdatedAt: p.UserUpdatedAt.Time,
			Email:     p.UserEmail,
			Username:  p.UserUsername,
			Password:  "",
			Role:      user.UserRole(p.UserRole),
		},
	}
}

func pgToApiPost(p *sqli.Post) *post.Post {
	topics := []post.PostTopic{}
	if err := json.Unmarshal(p.Topics, &topics); err != nil {
		slog.Error("Failed to unmarshal post topics", "error", err)
		topics = []post.PostTopic{}
	}

	return &post.Post{
		ID:          p.ID.Bytes,
		UserId:      p.UserID.Bytes,
		ThumbImage:  apiuuid(p.ThumbImage),
		CreatedAt:   p.CreatedAt.Time,
		UpdatedAt:   p.UpdatedAt.Time,
		Title:       p.Title,
		Description: p.Description,
		Content:     utils.B2S(p.Content),
		Topics:      topics,
	}
}

func pguuid(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func apiuuid(id pgtype.UUID) uuid.NullUUID {
	newId := uuid.NullUUID{Valid: false, UUID: uuid.Nil}
	if id.Valid {
		newId.Valid = true
		newId.UUID = id.Bytes
	}

	return newId
}

func str(s string) *string {
	return &s
}
