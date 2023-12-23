package repository

import (
	"context"
	"final-project/internal/driver/model"
)

type Driver interface {
	GetTripsList(ctx context.Context, userId string) (*[]model.Trip, error)
	GetTrip(ctx context.Context, userId string, tripId string) (*model.Trip, error)
	ChangeTripStatus(ctx context.Context, userId string, tripId string, newStatus string, cancelReason *string) error
}
