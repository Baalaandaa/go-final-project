package service

import (
	"context"
	"final-project/internal/driver/model"
)

type Driver interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
	ListTrips(ctx context.Context, userId string) (*[]model.Trip, error)
	GetTrip(ctx context.Context, userId string, tripId string) (*model.Trip, error)
	CancelTrip(ctx context.Context, userId string, tripId string, reason *string) error
	AcceptTrip(ctx context.Context, userId string, tripId string) error
	StartTrip(ctx context.Context, userId string, tripId string) error
	EndTrip(ctx context.Context, userId string, tripId string) error
}
