package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Yuruka00/go-user-subs/internal/handler"
	postgres_repo "github.com/Yuruka00/go-user-subs/internal/repository/postgres"
	"github.com/Yuruka00/go-user-subs/internal/service"
	"github.com/Yuruka00/go-user-subs/internal/tools/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	loggerHandler := slog.NewJSONHandler(os.Stdout, nil)
	baseLogger := slog.New(loggerHandler)

	cfg, err := config.Load()
	if err != nil {
		baseLogger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		baseLogger.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}

	repo := postgres_repo.NewSubscriptionRepository(db, baseLogger.With("layer", "repository"))
	srv := service.NewSubscriptionService(repo, baseLogger.With("layer", "service"))
	_ = handler.NewSubscriptionHandler(srv, baseLogger.With("layer", "handler"))

	err = http.ListenAndServe(":"+cfg.AppPort, nil)
	if err != nil {
		baseLogger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
