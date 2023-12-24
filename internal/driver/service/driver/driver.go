package driver

import (
	"context"
	"final-project/internal/driver/model"
	"final-project/internal/driver/repository"
	"final-project/internal/driver/service"
	kafka_producer "final-project/pkg/kafka-producer"
)

const (
	StatusDriverSearch = "DRIVER_SEARCH"
	StatusDriverFound  = "DRIVER_FOUND"
	StatusOnPosition   = "ON_POSITION"
	StatusStarted      = "STARTED"
	StatusEnded        = "ENDED"
	StatusCanceled     = "CANCELED"
)

type driverService struct {
	producer kafka_producer.Producer
	repo     repository.Driver
}

func (d driverService) CreateTrip(ctx context.Context, trip *model.Trip) error {
	return d.repo.CreateTrip(ctx, trip)
}

func (d driverService) ListTrips(ctx context.Context, userId string) (*[]model.Trip, error) {
	return d.repo.GetTripsList(ctx, userId)
}

func (d driverService) GetTrip(ctx context.Context, userId string, tripId string) (*model.Trip, error) {
	return d.repo.GetTrip(ctx, userId, tripId)
}

func (d driverService) CancelTrip(ctx context.Context, userId string, tripId string, reason *string) error {
	return d.repo.ChangeTripStatus(ctx, userId, tripId, StatusCanceled, reason)
}

func (d driverService) AcceptTrip(ctx context.Context, userId string, tripId string) error {
	return d.repo.ChangeTripStatus(ctx, userId, tripId, StatusDriverFound, nil)
}

func (d driverService) StartTrip(ctx context.Context, userId string, tripId string) error {
	return d.repo.ChangeTripStatus(ctx, userId, tripId, StatusStarted, nil)
}

func (d driverService) EndTrip(ctx context.Context, userId string, tripId string) error {
	return d.repo.ChangeTripStatus(ctx, userId, tripId, StatusEnded, nil)
}

func New(repo repository.Driver, producer kafka_producer.Producer) service.Driver {
	return &driverService{
		repo:     repo,
		producer: producer,
	}
}
