package cmd

import (
	"log/slog"
	"os"

	"github.com/zanz1n/go-htmx/internal/server"
)

func Run() {
	s := server.NewServer()
	addr := os.Getenv("LISTEN_ADDR")

	if err := s.Listen(addr); err != nil {
		slog.Error("Failed to listen on address: " + addr)
	}
}
