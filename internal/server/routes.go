package server

import "github.com/gofiber/fiber/v2"

func (s *Server) HandleHome(c *fiber.Ctx) error {
	return c.Render("index", struct{ Title string }{
		Title: "Home",
	})
}
