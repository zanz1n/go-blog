package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/fiberutils"
)

func (s *Server) HandleGetPost(c *fiber.Ctx) error {
	id := c.Params("id")

	p, err := s.postHandlers.HandleGetById(id)
	if err != nil {
		e, ok := err.(*errors.StatusError)
		if !ok {
			e = errors.ErrPostFetchFailed
		}

		if e.Code == errors.ErrPostFetchFailed.Code {
			return c.Status(500).Render("500",
				fiberutils.CreateProps(s.pp, c, "500", 0),
			)
		}

		return c.Status(404).Render("404",
			fiberutils.CreateProps(s.pp, c, "Post not find", 0),
		)
	}

	return c.Status(404).Render("post",
		fiberutils.CreateProps(s.pp, c, p.Title, p),
	)
}
