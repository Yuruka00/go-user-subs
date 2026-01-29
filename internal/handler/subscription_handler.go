package handler

import (
	"log/slog"

	"github.com/Yuruka00/go-user-subs/internal/service"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
	lg  *slog.Logger
}

func NewSubscriptionHandler(s *service.SubscriptionService, l *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		svc: s,
		lg:  l,
	}
}
