package pages

import "time"

type CreateRouteInfo struct {
	Name string
	Href string
}

type Route struct {
	IsCurrent bool
	Name      string
	Href      string
}

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
	Routes   []Route
	Data     T
}
