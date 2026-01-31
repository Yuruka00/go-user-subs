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
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Logger Initialization
	loggerHandler := slog.NewJSONHandler(os.Stdout, nil)
	baseLogger := slog.New(loggerHandler)

	// Config Reading
	cfg, err := config.Load()
	if err != nil {
		baseLogger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

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
	_ = handler.NewSubscriptionHandler(srv, baseLogger.With("layer", "handler"))

	// Server starting
	err = http.ListenAndServe(":"+cfg.AppPort, nil)
	if err != nil {
		baseLogger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
