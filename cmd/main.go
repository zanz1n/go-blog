package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/zanz1n/go-htmx/internal/server"
)

func Run() {
	sysch := make(chan os.Signal, 1)
	signal.Notify(sysch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	s := server.NewServer(os.Getenv("APP_NAME"))
	addr := os.Getenv("LISTEN_ADDR")

	go func() {
		if err := s.Listen(addr); err != nil {
			slog.Error("Failed to listen on address: " + addr)
		}
	}()

	<-sysch
	slog.Info("Shutting down ...")
	s.Shutdown()
}
