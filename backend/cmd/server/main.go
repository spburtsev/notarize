package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"

	"github.com/spburtsev/notarize/internal/auth"
	"github.com/spburtsev/notarize/internal/config"
	"github.com/spburtsev/notarize/internal/db"
	"github.com/spburtsev/notarize/internal/db/models"
	"github.com/spburtsev/notarize/internal/handler"
	"github.com/spburtsev/notarize/internal/oas"
	"github.com/spburtsev/notarize/internal/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))
	slog.SetDefault(logger)

	gdb, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect to database", "error", err)
		os.Exit(1)
	}
	if sqlDB, err := gdb.DB(); err == nil {
		defer sqlDB.Close()
	}

	if err := gdb.AutoMigrate(models.All...); err != nil {
		slog.Error("run migrations", "error", err)
		os.Exit(1)
	}

	if err := seedAdmin(ctx, gdb, cfg); err != nil {
		slog.Error("seed admin", "error", err)
		os.Exit(1)
	}

	store, err := storage.New(storage.Config{
		Endpoint:  cfg.S3Endpoint,
		AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey,
		Bucket:    cfg.S3Bucket,
		UseSSL:    cfg.S3UseSSL,
	})
	if err != nil {
		slog.Error("connect to object storage", "error", err)
		os.Exit(1)
	}
	if err := store.EnsureBucket(ctx); err != nil {
		slog.Error("ensure bucket", "error", err)
		os.Exit(1)
	}

	authSvc := auth.NewService(cfg.JWTSecret, cfg.JWTTTL)
	service := handler.New(gdb, store, authSvc)
	oasHandler, err := oas.NewServer(service, authSvc, oas.WithPathPrefix("/api/v1"))
	if err != nil {
		slog.Error("create server", "error", err)
		os.Exit(1)
	}

	srv := &http.Server{Addr: ":8080", Handler: oasHandler}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		stop()
		slog.Info("shutdown signal received, draining connections")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
		}
		slog.Info("server stopped")
	}
}

func seedAdmin(ctx context.Context, gdb *gorm.DB, cfg config.Config) error {
	if cfg.SeedAdminEmail == "" || cfg.SeedAdminPassword == "" {
		return nil
	}

	var count int64
	if err := gdb.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", cfg.SeedAdminEmail).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := auth.Hash(cfg.SeedAdminPassword)
	if err != nil {
		return err
	}
	admin := models.User{
		Email:        cfg.SeedAdminEmail,
		FirstName:    "Admin",
		LastName:     "User",
		PasswordHash: hash,
		Role:         models.UserRoleAdmin,
	}
	if err := gdb.WithContext(ctx).Create(&admin).Error; err != nil {
		return err
	}
	slog.Info("seeded admin user", "email", cfg.SeedAdminEmail)
	return nil
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
