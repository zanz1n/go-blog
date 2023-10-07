package handlers

import (
	"github.com/google/uuid"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/post"
	"github.com/zanz1n/go-htmx/internal/post/repository"
)

func NewPostHandlers(ps repository.PostRepository) *PostHandlers {
	return &PostHandlers{
		ps: ps,
	}
}

type PostHandlers struct {
	ps repository.PostRepository
}

func (h *PostHandlers) HandleGetById(id string) (*post.PostWithUser, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.ErrInvalidUUID
	}

	return h.ps.GetByIdWithUser(uid)
}
