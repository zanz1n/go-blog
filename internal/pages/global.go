package pages

import "time"

type UserProps struct {
	ID         string
	Username   string
	Email      string
	CreatedAt  time.Time
	Expiration time.Time
}

type PageProps[T any] struct {
	AppName  string
	PageName string
	User     *UserProps // Nullable
	Data     T
}
