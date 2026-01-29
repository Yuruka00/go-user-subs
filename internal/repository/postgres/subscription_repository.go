package postgres

import (
	"context"
	"log/slog"

	"github.com/Yuruka00/go-user-subs/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
	lg *slog.Logger
}

func NewSubscriptionRepository(d *gorm.DB, l *slog.Logger) *subscriptionRepository {
	return &subscriptionRepository{
		db: d,
		lg: l,
	}
}

func (sr *subscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	return nil
}
func (sr *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return nil, nil
}
func (sr *subscriptionRepository) Update(ctx context.Context, s *domain.Subscription) error {
	return nil
}
func (sr *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (sr *subscriptionRepository) GetList(ctx context.Context, f *domain.SubscriptionFilter) ([]*domain.Subscription, error) {
	return nil, nil
}

func (sr *subscriptionRepository) SumByFilter(ctx context.Context, f *domain.SubscriptionFilter) (int, error) {
	return 0, nil
}
