package cmd

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/zanz1n/go-htmx/internal/fiberutils"
	"github.com/zanz1n/go-htmx/website"
)

var (
	assets    = http.FS(website.EmbedAssets)
	assetsDir = "/dist/templates"
)

func Run() {
	engine := handlebars.NewFileSystem(assets, ".hbs")
	engine.Directory = assetsDir

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

	app.Use(fiberutils.NewLoggerMiddleware())

	app.Hooks().OnListen(func(ld fiber.ListenData) error {
		slog.Info("Listenning for requests", "port", ld.Port)
		return nil
	})

	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       assets,
		Browse:     false,
		PathPrefix: "dist",
		MaxAge:     3600,
		Next: func(c *fiber.Ctx) bool {
			if strings.HasPrefix(c.Path(), "/assets/templates") {
				return true
			}

			return false
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Home",
		})
	})

	app.Listen(os.Getenv("LISTEN_ADDR"))
}
