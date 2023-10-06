package repository

import (
	"github.com/google/uuid"
	"github.com/zanz1n/go-htmx/internal/post"
)

type PostRepository interface {
	GetById(id uuid.UUID) (*post.Post, error)
	GetByIdWithUser(id uuid.UUID) (*post.PostWithUser, error)
	GetPartialById(id uuid.UUID) (*post.PartialPost, error)
	Create(userId uuid.UUID, data *post.PostCreateData) (*post.Post, error)
}
