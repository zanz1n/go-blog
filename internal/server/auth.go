package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/go-htmx/internal/auth/handlers"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/fiberutils"
	"github.com/zanz1n/go-htmx/internal/pages"
	"github.com/zanz1n/go-htmx/internal/user"
)

func (s *Server) UserExtractorMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("auth_token")
	if cookie == "" {
		return c.Next()
	}

	u, err := s.authHandlers.DecodeAuthToken(cookie)
	if err == nil {
		c.Locals("user", u)
	}

	return c.Next()
}

func (s *Server) HandleGetHome(c *fiber.Ctx) error {
	return c.Status(200).Render("index",
		fiberutils.CreateProps(s.pp, c, "Home", 0),
	)
}

func (s *Server) HandleGetLogin(c *fiber.Ctx) error {
	return c.Status(200).Render("login",
		fiberutils.CreateProps(s.pp, c, "Login", pages.LoginPageData{
			Error: nil,
		}),
	)
}

func (s *Server) HandlePostLogin(c *fiber.Ctx) error {
	data := handlers.LoginIdenPayload{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(401).Render("login",
			fiberutils.CreateProps(s.pp, c, "Login", pages.LoginPageData{
				Error: str("Invalid form data"),
			}),
		)
	}

	token, err := s.authHandlers.HandleLogin(&data)
	if err != nil {
		e, ok := err.(*errors.StatusError)
		if !ok {
			e = errors.ErrInternalServerError
		}

		return c.Status(e.HttpCode).Render("login",
			fiberutils.CreateProps(s.pp, c, "Login", pages.LoginPageData{
				Error: str(fmt.Sprintf("Error %d: %s", e.Code, e.Message)),
			}),
		)
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token

	c.Cookie(cookie)

	return c.Redirect("/" + c.Query("from"))
}

func (s *Server) HandleGetSignup(c *fiber.Ctx) error {
	return c.Status(200).Render("signup",
		fiberutils.CreateProps(s.pp, c, "Sign Up", 0),
	)
}

func (s *Server) HandlePostSignup(c *fiber.Ctx) error {
	data := user.UserCreateData{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(401).Render("signup",
			fiberutils.CreateProps(s.pp, c, "Sign Up", pages.SignupPageData{
				Error: str("Invalid form data"),
			}),
		)
	}

	_, token, err := s.authHandlers.HandleSignup(&data)
	if err != nil {
		e, ok := err.(*errors.StatusError)
		if !ok {
			e = errors.ErrInternalServerError
		}

		return c.Status(e.HttpCode).Render("signup",
			fiberutils.CreateProps(s.pp, c, "Sign Up", pages.SignupPageData{
				Error: str(fmt.Sprintf("Error %d: %s", e.Code, e.Message)),
			}),
		)
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token

	c.Cookie(cookie)

	return c.Redirect("/" + c.Query("from"))
}

func (s *Server) HandleLogout(c *fiber.Ctx) error {
	c.ClearCookie("auth_token")
	return c.Redirect("/" + c.Query("from"))
}
