package service

import (
	"context"
	"log/slog"

	"github.com/Yuruka00/go-user-subs/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetList(ctx context.Context, f *domain.SubscriptionFilter) ([]domain.Subscription, error)
	SumByFilter(ctx context.Context, f *domain.SubscriptionFilter) (int, error)
}

type SubscriptionService struct {
	repo SubscriptionRepository
	lg   *slog.Logger
}

func NewSubscriptionService(r SubscriptionRepository, l *slog.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo: r,
		lg:   l,
	}
}
