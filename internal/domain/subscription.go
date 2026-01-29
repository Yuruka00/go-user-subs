package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ServiceName string
	UserID      uuid.UUID
	Price       int
	DateStart   time.Time
	DateEnd     *time.Time
}

type SubscriptionFilter struct {
	ServiceName *string
	UserID      *uuid.UUID
	DateFrom    time.Time
	DateTo      time.Time
}
