package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/spburtsev/notarize/internal/config"
	"github.com/spburtsev/notarize/internal/db"
	"github.com/spburtsev/notarize/internal/handler"
	"github.com/spburtsev/notarize/internal/oas"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))
	slog.SetDefault(logger)

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.Migrate(pool); err != nil {
		slog.Error("run migrations", "error", err)
		os.Exit(1)
	}

	service := &handler.ServerHandler{}
	srv, err := oas.NewServer(service, nil)
	if err != nil {
		slog.Error("create server", "error", err)
		os.Exit(1)
	}

	slog.Info("starting server", "addr", ":8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func parseLogLevel(s string) slog.Level {
	if s == "" {
		return slog.LevelInfo
	}
	var level slog.Level
	if err := level.UnmarshalText([]byte(s)); err != nil {
		return slog.LevelInfo
	}
	return level
}
