package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	UserID      uuid.UUID
	price       int
	dateStart   time.Time
	dateEnd     time.Time
}

type SubscriptionFilter struct {
	ServiceName string
	UserID      uuid.UUID
	dateFrom    time.Time
	dateTo      time.Time
}
