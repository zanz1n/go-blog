package fiberutils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/go-htmx/internal/pages"
)

func CreateProps[T any](
	pp *pages.PagePropsProvider,
	c *fiber.Ctx,
	data T,
) pages.PageProps[T] {
	pageName := ""

	routes := make([]pages.Route, len(pp.Routes))
	for i, v := range pp.Routes {
		isCurrent := strings.Contains(c.Path(), v.Href)
		if isCurrent {
			pageName = v.Name
		}

		routes[i] = pages.Route{
			IsCurrent: isCurrent,
			Name:      v.Name,
			Href:      v.Href,
		}
	}

	return pages.PageProps[T]{
		AppName:  pp.AppName,
		PageName: pageName,
		Routes:   routes,
		User:     nil,
		Data:     data,
	}
}
