package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/zanz1n/go-htmx/internal/pages"
	"github.com/zanz1n/go-htmx/internal/server"
)

func Run() {
	sysch := make(chan os.Signal, 1)
	signal.Notify(sysch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	pp := pages.PagePropsProvider{
		AppName: os.Getenv("APP_NAME"),
	}

	s := server.NewServer(&pp)
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
