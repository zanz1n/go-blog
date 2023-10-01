package fiberutils

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		chainErr := c.Next()

		errHandler := c.App().ErrorHandler

		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				slog.Error("Not Handled Fiber Error", "error", chainErr)
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		slog.Info("Http Request",
			"addr", c.Context().RemoteAddr(),
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration", time.Since(start),
		)

		return nil
	}
}
