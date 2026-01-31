package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yuruka00/go-user-subs/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(d *gorm.DB) *subscriptionRepository {
	return &subscriptionRepository{
		db: d,
	}
}

func (sr *subscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	err := gorm.G[domain.Subscription](sr.db).Create(ctx, s)
	if err != nil {
		return fmt.Errorf("subscriptionRepository.Create: %w", err)
	}
	return nil
}

func (sr *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	s, err := sr.getByIDQuery(id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subscriptionRepository.GetByID: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("subscriptionRepository.GetByID: %w", err)
	}
	return &s, nil
}
func (sr *subscriptionRepository) Update(ctx context.Context, s *domain.Subscription) error {
	ra, err := sr.getByIDQuery(s.ID).Omit("id").Updates(ctx, *s)
	if err != nil {
		return fmt.Errorf("subscriptionRepository.Update: %w", err)
	}
	if ra == 0 {
		return fmt.Errorf("subscriptionRepository.Update: %w", domain.ErrNotFound)
	}
	return nil
}
func (sr *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ra, err := sr.getByIDQuery(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("subscriptionRepository.Delete: %w", err)
	}
	if ra == 0 {
		return fmt.Errorf("subscriptionRepository.Delete: %w", domain.ErrNotFound)
	}
	return nil
}
func (sr *subscriptionRepository) GetList(ctx context.Context, f *domain.SubscriptionFilter) ([]domain.Subscription, error) {
	rs, err := sr.getFilteredQuery(f).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("subscriptionRepository.GetList: %w", err)
	}
	return rs, nil
}

func (sr *subscriptionRepository) SumByFilter(ctx context.Context, f *domain.SubscriptionFilter) (int, error) {
	var total int
	err := sr.getFilteredQuery(f).Select("COALESCE(SUM(price), 0)").Scan(ctx, &total)
	if err != nil {
		return 0, fmt.Errorf("subscriptionRepository.SumByFilter: %w", err)
	}
	return total, nil
}

func (sr *subscriptionRepository) getByIDQuery(id uuid.UUID) gorm.ChainInterface[domain.Subscription] {
	return gorm.G[domain.Subscription](sr.db).Where("id = ?", id)
}

func (sr *subscriptionRepository) getFilteredQuery(f *domain.SubscriptionFilter) gorm.ChainInterface[domain.Subscription] {
	query := gorm.G[domain.Subscription](sr.db).Where("1=1")

	if f.ServiceName != nil {
		query = query.Where("service_name = ?", *f.ServiceName)
	}
	if f.UserID != nil {
		query = query.Where("user_id = ?", *f.UserID)
	}

	return query.Where("date_start BETWEEN ? AND ?", f.DateFrom, f.DateTo)
}
