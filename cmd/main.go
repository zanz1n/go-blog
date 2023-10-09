package cmd

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/zanz1n/go-htmx/internal/auth"
	auth_handlers "github.com/zanz1n/go-htmx/internal/auth/handlers"
	post_handlers "github.com/zanz1n/go-htmx/internal/post/handlers"
	post_repository "github.com/zanz1n/go-htmx/internal/post/repository"
	"github.com/zanz1n/go-htmx/internal/server"
	"github.com/zanz1n/go-htmx/internal/sqli"
	user_repository "github.com/zanz1n/go-htmx/internal/user/repository"
)

func Run() {
	sysch := make(chan os.Signal, 1)
	signal.Notify(sysch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connStart := time.Now()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	slog.Info("Connected to postgres", "duration", time.Since(connStart))

	dba := sqli.New(conn)
	htmlParser := post_repository.NewHtmlParser(true, true)

	userRepository := user_repository.NewPostgresRepository(dba)
	authRepository := auth.NewJwtRepository([]byte(os.Getenv("JWT_HMAC_KEY")), time.Hour)
	postRepository := post_repository.NewPostgresRepository(dba, htmlParser)

	s := server.NewServer(
		os.Getenv("APP_NAME"),
		auth_handlers.NewAuthHandlers(authRepository, userRepository),
		post_handlers.NewPostHandlers(postRepository),
	)
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
