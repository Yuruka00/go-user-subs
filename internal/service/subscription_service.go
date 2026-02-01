package service

import (
	"context"
	"fmt"
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

func (ss *SubscriptionService) Create(ctx context.Context, s *domain.Subscription) error {
	ss.lg.Debug("creating subscription",
		"data", s,
	)

	err := ss.repo.Create(ctx, s)

	if err != nil {
		return fmt.Errorf("SubscriptionService.Create(): %w", err)
	}
	return nil
}
func (ss *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	s, err := ss.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("SubscriptionService.GetByID(id: %s): %w", id, err)
	}
	return s, nil
}
func (ss *SubscriptionService) Update(ctx context.Context, s *domain.Subscription) error {
	ss.lg.Debug("updating subscription",
		"data", s,
	)

	err := ss.repo.Update(ctx, s)

	if err != nil {
		return fmt.Errorf("SubscriptionService.Update(id: %s): %w", s.ID, err)
	}
	return nil
}
func (ss *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	err := ss.repo.Delete(ctx, id)

	if err != nil {
		return fmt.Errorf("SubscriptionService.Delete(id: %s): %w", id, err)
	}
	return nil
}
func (ss *SubscriptionService) GetList(ctx context.Context, sf *domain.SubscriptionFilter) ([]domain.Subscription, error) {
	ss.lg.Debug("getting subscription list",
		"filter", sf,
	)

	rs, err := ss.repo.GetList(ctx, sf)

	if err != nil {
		return nil, fmt.Errorf("SubscriptionService.GetList: %w", err)
	}
	return rs, nil
}
func (ss *SubscriptionService) CalculateTotalPrice(ctx context.Context, sf *domain.SubscriptionFilter) (int, error) {
	ss.lg.Debug("calculating subscription total price",
		"filter", sf,
	)

	rs, err := ss.repo.SumByFilter(ctx, sf)

	if err != nil {
		return 0, fmt.Errorf("SubscriptionService.CalculateTotalPrice: %w", err)
	}
	return rs, nil
}
