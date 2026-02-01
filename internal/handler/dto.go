package handler

import (
	"time"

	"github.com/Yuruka00/go-user-subs/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionCreateRequest struct {
	ServiceName string    `json:"service_name" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Price       int       `json:"price" validate:"gte=0"`
	DateStart   string    `json:"date_start" validate:"required,datetime=2006-01"`
	DateEnd     *string   `json:"date_end" validate:"omitempty,datetime=2006-01-02"`
}

func (scr *SubscriptionCreateRequest) ToDomain() (*domain.Subscription, error) {
	dStart, err := time.Parse("2006-01", scr.DateStart)
	if err != nil {
		return nil, err
	}

	var dEnd *time.Time
	if scr.DateEnd != nil {
		tmp, err := time.Parse("2006-01-02", *scr.DateEnd)
		if err != nil {
			return nil, err
		}
		dEnd = &tmp
	}

	item := domain.Subscription{
		ServiceName: scr.ServiceName,
		UserID:      scr.UserID,
		Price:       scr.Price,
		DateStart:   dStart,
		DateEnd:     dEnd,
	}
	return &item, nil
}

type SubscriptionUpdateRequest struct {
	ServiceName string    `json:"service_name" validate:"omitempty"`
	UserID      uuid.UUID `json:"user_id" validate:"omitempty"`
	Price       *int      `json:"price" validate:"omitempty,gte=0"`
	DateStart   string    `json:"date_start" validate:"omitempty,datetime=2006-01"`
	DateEnd     string    `json:"date_end" validate:"omitempty,datetime=2006-01-02"`
}

func (scr *SubscriptionUpdateRequest) ToDomain() (*domain.Subscription, bool, error) {
	var item domain.Subscription
	hasChanges := false

	if scr.ServiceName != "" {
		item.ServiceName = scr.ServiceName
		hasChanges = true
	}

	if scr.UserID != uuid.Nil {
		item.UserID = scr.UserID
		hasChanges = true
	}

	if scr.Price != nil {
		item.Price = *scr.Price
		hasChanges = true
	}

	if scr.DateStart != "" {
		dStart, err := time.Parse("2006-01", scr.DateStart)
		if err != nil {
			return nil, false, err
		}
		item.DateStart = dStart
		hasChanges = true
	}

	if scr.DateEnd != "" {
		dEnd, err := time.Parse("2006-01", scr.DateEnd)
		if err != nil {
			return nil, false, err
		}
		item.DateEnd = &dEnd
		hasChanges = true
	}

	return &item, hasChanges, nil
}

type SubscriptionFilterRequest struct {
	ServiceName string    `json:"service_name"`
	UserID      uuid.UUID `json:"user_id"`
	DateFrom    string    `json:"date_from"`
	DateTo      string    `json:"date_to"`
}

func (sfr *SubscriptionFilterRequest) ToDomain() (*domain.SubscriptionFilter, error) {
	var item domain.SubscriptionFilter

	if sfr.ServiceName != "" {
		item.ServiceName = &sfr.ServiceName
	}
	if sfr.UserID != uuid.Nil {
		item.UserID = &sfr.UserID
	}

	dFrom, err := time.Parse("2006-01", sfr.DateFrom)
	if err != nil {
		return nil, err
	}
	item.DateFrom = dFrom

	dTo, err := time.Parse("2006-01", sfr.DateTo)
	if err != nil {
		return nil, err
	}
	item.DateTo = dTo

	return &item, nil
}
