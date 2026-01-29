package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Yuruka00/go-user-subs/internal/tools/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	loggerHandler := slog.NewJSONHandler(os.Stdout, nil)
	l := slog.New(loggerHandler)

	cfg, err := config.Load()
	if err != nil {
		l.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	_, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		l.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}

	err = http.ListenAndServe(":"+cfg.AppPort, nil)
	if err != nil {
		l.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
