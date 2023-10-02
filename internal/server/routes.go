package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/go-htmx/internal/fiberutils"
)

func (s *Server) HandleHome(c *fiber.Ctx) error {
	return c.Status(200).Render("index",
		fiberutils.CreateProps(s.pp, c, "Hello World"),
	)
}

func (s *Server) HandleLogin(c *fiber.Ctx) error {
	return c.Status(200).Render("login",
		fiberutils.CreateProps(s.pp, c, "Hello World"),
	)
}
