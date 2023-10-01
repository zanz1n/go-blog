package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/fiberutils"
	"github.com/zanz1n/go-htmx/website"
)

type Server struct {
	app *fiber.App
}

func NewServer() *Server {
	fs := http.FS(website.EmbedAssets)
	s := Server{}

	engine := handlebars.NewFileSystem(fs, ".hbs")
	engine.Directory = "/dist/templates"

	app := fiber.New(fiber.Config{
		ServerHeader:          "fasthttp",
		CaseSensitive:         true,
		Prefork:               false,
		StrictRouting:         false,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		Views:                 engine,
	})
	s.app = app

	app.Hooks().OnListen(s.OnListen)

	app.Use(fiberutils.NewLoggerMiddleware())
	app.Use("/assets", s.assetsHandler(fs))

	app.Get("/", s.HandleHome)

	return &s
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
