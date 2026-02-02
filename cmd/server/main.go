package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Yuruka00/go-user-subs/internal/handler"
	postgres_repo "github.com/Yuruka00/go-user-subs/internal/repository/postgres"
	"github.com/Yuruka00/go-user-subs/internal/service"
	"github.com/Yuruka00/go-user-subs/internal/tools/config"
	"github.com/Yuruka00/go-user-subs/migrations"
	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Logger Initialization
	baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Config Reading
	cfg, err := config.Load()
	if err != nil {
		baseLogger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Logging level setup
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		logLevel = slog.LevelInfo
	}
	baseLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	// Database Conection Establishing
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		baseLogger.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}

	sqldb, err := db.DB()
	if err != nil {
		baseLogger.Error("failed to get *sql.DB object", "error", err)
		os.Exit(1)
	}

	// Goose Migrations
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		baseLogger.Error("failed to set migrations dialect", "error", err)
		os.Exit(1)
	}

	if err := goose.Up(sqldb, "."); err != nil {
		baseLogger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Dependency Injection
	repo := postgres_repo.NewSubscriptionRepository(db)
	srv := service.NewSubscriptionService(repo, baseLogger.With("layer", "service"))
	h := handler.NewSubscriptionHandler(srv, baseLogger.With("layer", "handler"))

	// Routes Setup
	r := chi.NewRouter()

	r.Post("/subscriptions", h.Create)
	r.Get("/subscriptions/{id}", h.Get)
	r.Patch("/subscriptions/{id}", h.Update)
	r.Delete("/subscriptions/{id}", h.Delete)
	r.Get("/subscriptions", h.GetList)
	r.Get("/subscriptions/total", h.GetTotalPrice)

	// Server starting
	err = http.ListenAndServe(":"+cfg.AppPort, r)
	if err != nil {
		baseLogger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
