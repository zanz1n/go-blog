package post

import (
	"time"

	"github.com/google/uuid"
	"github.com/zanz1n/go-htmx/internal/user"
)

type PostTopicKind string

const (
	PostTopicKindH2 PostTopicKind = "h2"
	PostTopicKindH3 PostTopicKind = "h3"
	PostTopicKindH4 PostTopicKind = "h4"
	PostTopicKindH5 PostTopicKind = "h5"
)

type PostTopic struct {
	Title string        `json:"title"`
	Kind  PostTopicKind `json:"kind"`
}

type Post struct {
	ID          uuid.UUID     `json:"id"`
	UserId      uuid.UUID     `json:"user_id"`
	ThumbImage  uuid.NullUUID `json:"thumb_image"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Content     string        `json:"content"`
	Topics      []PostTopic   `json:"topics"`
}

type PartialPost struct {
	ID           uuid.UUID     `json:"id"`
	UserId       uuid.UUID     `json:"user_id"`
	UserUsername string        `json:"user_username"`
	ThumbImage   uuid.NullUUID `json:"thumb_image"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
}

type PostWithUser struct {
	Post
	User user.User `json:"user"`
}

type PostCreateData struct {
	Title       string    `json:"title" validate:"required"`
	Description *string   `json:"description"`
	Content     string    `json:"content" validate:"required"`
	ThumbImage  uuid.UUID `json:"thumb_image" validate:"required"`
}
