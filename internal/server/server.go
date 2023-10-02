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
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/fiberutils"
	"github.com/zanz1n/go-htmx/internal/pages"
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
	app *fiber.App
	pp  *pages.PagePropsProvider
}

func NewServer(appName string) *Server {
	fs := http.FS(website.EmbedAssets)

	s := Server{
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

	s.app.Use(recover.New())
	s.app.Use(fiberutils.NewLoggerMiddleware())
	s.app.Use("/assets", s.assetsHandler(fs))

	s.app.Get("/", s.HandleHome)

	return &s
}

func (s *Server) ErrorHandler(c *fiber.Ctx, err error) error {
	mt, _, _ := mime.ParseMediaType(c.Accepts())

	e, ok := err.(*fiber.Error)
	if !ok {
		slog.Error("Unhandled error", "error", err)
		e = fiber.ErrInternalServerError
	}
	s.handleErrorJson(c, e)

	if mt == "application/json" {
		s.handleErrorJson(c, e)
	} else {
		s.handleHtmlError(c, e)
	}

	return nil
}

func (s *Server) handleHtmlError(c *fiber.Ctx, e *fiber.Error) {
	if e.Code == 404 {
		c.Status(404).Render("404",
			fiberutils.CreateProps(s.pp, c, ""),
		)
	} else {
		c.Status(500).Render("500",
			fiberutils.CreateProps(s.pp, c, ""),
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.app.ShutdownWithContext(ctx)
}

func (s *Server) OnListen(ld fiber.ListenData) error {
	slog.Info("Listenning for requests", "port", ld.Port)
	return nil
}

func (s *Server) HandleErr(c *fiber.Ctx, err error) error {
	e := errors.GetStatusErr(err)

	return c.Status(e.HttpCode()).JSON(errors.ErrorBody{
		Message:   e.Message(),
		ErrorCode: e.CustomCode(),
	})
}
