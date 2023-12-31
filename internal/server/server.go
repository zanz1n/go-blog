package server

import (
	"context"
	"log/slog"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/handlebars/v2"
	auth_handlers "github.com/zanz1n/go-htmx/internal/auth/handlers"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/fiberutils"
	"github.com/zanz1n/go-htmx/internal/pages"
	post_handlers "github.com/zanz1n/go-htmx/internal/post/handlers"
	"github.com/zanz1n/go-htmx/website"
)

var routes = []pages.CreateRouteInfo{
	{
		Name: "Home",
		Href: "/",
	},
	{
		Name: "Posts",
		Href: "/posts",
	},
}

type Server struct {
	app          *fiber.App
	pp           *pages.PagePropsProvider
	authHandlers *auth_handlers.AuthHandlers
	postHandlers *post_handlers.PostHandlers
}

func NewServer(
	appName string,
	ah *auth_handlers.AuthHandlers,
	ph *post_handlers.PostHandlers,
) *Server {
	fs := http.FS(website.EmbedAssets)

	s := Server{
		authHandlers: ah,
		postHandlers: ph,
		pp: &pages.PagePropsProvider{
			AppName: appName,
			Routes:  routes,
		},
	}

	engine := handlebars.NewFileSystem(fs, ".hbs")
	engine.Directory = "/dist/templates"

	s.app = fiber.New(fiber.Config{
		ServerHeader:          "fasthttp",
		CaseSensitive:         true,
		Prefork:               false,
		StrictRouting:         false,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		Views:                 engine,
		AppName:               appName,
		ErrorHandler:          s.ErrorHandler,
	})

	s.app.Hooks().OnListen(s.OnListen)

	// s.app.Use(recover.New())
	s.app.Use(fiberutils.NewLoggerMiddleware())
	s.app.Use(cors.New())
	s.app.Use("/assets", s.assetsHandler(fs))

	s.app.Use(s.UserExtractorMiddleware)

	s.app.Get("/", s.HandleGetHome)

	s.app.Get("/login", s.HandleGetLogin)
	s.app.Post("/login", s.HandlePostLogin)
	s.app.Get("/signup", s.HandleGetSignup)
	s.app.Post("/signup", s.HandlePostSignup)
	s.app.Get("/logout", s.HandleLogout)
	s.app.Post("/logout", s.HandleLogout)

	s.app.Get("/post/:id", s.HandleGetPost)

	return &s
}

func (s *Server) ErrorHandler(c *fiber.Ctx, err error) error {
	mt, _, _ := mime.ParseMediaType(c.Accepts())

	e, ok := err.(*fiber.Error)
	if !ok {
		slog.Error("Unhandled error", "error", err)
		e = fiber.ErrInternalServerError
	}

	if mt == "application/json" {
		s.handleErrorJson(c, e)
	} else {
		s.handleHtmlError(c, e)
	}

	return nil
}

func (s *Server) handleHtmlError(c *fiber.Ctx, e *fiber.Error) {
	if e.Code == 403 || e.Code == 404 || e.Code == 405 {
		c.Status(404).Render("404",
			fiberutils.CreateProps(s.pp, c, "404", struct{}{}),
		)
	} else {
		c.Status(500).Render("500",
			fiberutils.CreateProps(s.pp, c, "500", struct{}{}),
		)
	}
}

func (s *Server) handleErrorJson(c *fiber.Ctx, e *fiber.Error) {
	c.Status(e.Code).JSON(errors.ErrorBody{
		Message:   e.Message,
		ErrorCode: uint(e.Code),
	})
}

func (s *Server) assetsHandler(fs http.FileSystem) func(*fiber.Ctx) error {
	return filesystem.New(filesystem.Config{
		Root:       fs,
		Browse:     false,
		PathPrefix: "dist",
		MaxAge:     3600,
		Next: func(c *fiber.Ctx) bool {
			if strings.HasPrefix(c.Path(), "/assets/templates") {
				s.ErrorHandler(c, fiber.ErrNotFound)
				return true
			}

			return false
		},
	})
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Fiber shutdown failed", "error", err)
	}
}

func (s *Server) OnListen(ld fiber.ListenData) error {
	slog.Info("Listenning for requests", "port", ld.Port)
	return nil
}

func str(s string) *string {
	return &s
}
