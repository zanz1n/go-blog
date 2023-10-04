package fiberutils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/go-htmx/internal/auth"
	"github.com/zanz1n/go-htmx/internal/pages"
)

func CreateProps[T any](
	pp *pages.PagePropsProvider,
	c *fiber.Ctx,
	pageName string,
	data T,
) pages.PageProps[T] {
	var user *pages.UserProps = nil

	if u, ok := c.Locals("user").(*auth.UserAuthPayload); ok {
		user = &pages.UserProps{
			ID:         u.UserId.String(),
			Email:      u.Email,
			Expiration: time.Unix(u.ExpiryDate, 0),
		}
	}

	return pages.CreateProps(pp, c.Path(), pageName, user, data)
}
