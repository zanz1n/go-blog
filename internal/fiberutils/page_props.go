package fiberutils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/go-htmx/internal/pages"
)

func CreateProps[T any](
	pp *pages.PagePropsProvider,
	c *fiber.Ctx,
	pageName string,
	data T,
) pages.PageProps[T] {
	return pages.CreateProps(pp, c.Path(), pageName, nil, data)
}
