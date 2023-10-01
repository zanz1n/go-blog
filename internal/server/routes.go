package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/go-htmx/internal/pages"
)

func (s *Server) HandleHome(c *fiber.Ctx) error {
	return c.Render("index", pages.CreateProps(
		s.pp,
		"Home",
		nil,
		"Hello World",
	))
}
