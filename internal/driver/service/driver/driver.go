package driver

import (
	"context"
	"final-project/internal/driver/model"
	"final-project/internal/driver/repository"
	"final-project/internal/driver/service"
	"final-project/pkg/helpers"
	kafka_producer "final-project/pkg/kafka-producer"
	location_client "final-project/pkg/location-client"
	"math"
)

const (
	StatusDriverSearch = "DRIVER_SEARCH"
	StatusDriverFound  = "DRIVER_FOUND"
	StatusOnPosition   = "ON_POSITION"
	StatusStarted      = "STARTED"
	StatusEnded        = "ENDED"
	StatusCanceled     = "CANCELED"
	MinRadius          = 0.01
	MaxRadius          = 2.0
	RadiusFactor       = 2.0
)

type driverService struct {
	producer        kafka_producer.Producer
	repo            repository.Driver
	locationBaseUrl string
}

func euclideanDistance(a, b model.LatLngLiteral) float64 {
	return math.Sqrt(math.Pow(b.Lat-a.Lat, 2) + math.Pow(b.Lng-a.Lng, 2))
}

func findClosestDriver(pos model.LatLngLiteral, drivers []location_client.Driver) *location_client.Driver {
	var closest *location_client.Driver
	minDistance := MaxRadius

	for _, driver := range drivers {
		distance := euclideanDistance(pos, model.LatLngLiteral{
			Lat: driver.Lat,
			Lng: driver.Lng,
		})
		if distance < minDistance {
			minDistance = distance
			closest = &driver
		}
	}

	return closest
}

func (d driverService) CreateTrip(ctx context.Context, trip *model.Trip) error {
	client := location_client.New(d.locationBaseUrl)

	radius := MinRadius
	var driver *location_client.Driver
	for radius < MaxRadius && driver == nil {
		drivers, err := client.GetDriverLocations(trip.From.Lat, trip.From.Lng, radius)
		if err != nil {
			return err
		}

		driver = findClosestDriver(trip.From, drivers)
		radius *= RadiusFactor
	}

	if driver == nil {
		trip.Status = StatusCanceled
		*trip.CancelReason = "driver not found"
	} else {
		trip.DriverId = driver.DriverId
	}

	return d.repo.CreateTrip(ctx, trip)
}

func (d driverService) ListTrips(ctx context.Context, userId string) ([]*model.Trip, error) {
	return d.repo.GetTripsList(ctx, userId)
}

func (d driverService) GetTrip(ctx context.Context, userId string, tripId string) (*model.Trip, error) {
	trip, err := d.repo.GetTrip(ctx, userId, tripId)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, helpers.ErrNotFound
	}
	return trip, nil
}

func (d driverService) CancelTrip(ctx context.Context, userId string, tripId string, reason *string) error {
	trip, err := d.repo.GetTrip(ctx, userId, tripId)
	if err != nil {
		return err
	}
	if trip == nil {
		return helpers.ErrNotFound
	}
	return d.repo.ChangeTripStatus(ctx, userId, tripId, StatusCanceled, reason)
}

func (d driverService) AcceptTrip(ctx context.Context, userId string, tripId string) error {
	trip, err := d.repo.GetTrip(ctx, userId, tripId)
	if err != nil {
		return err
	}
	if trip == nil {
		return helpers.ErrNotFound
	}
	err = d.repo.ChangeTripStatus(ctx, userId, tripId, StatusDriverFound, nil)
	if err != nil {
		return err
	}
	return d.producer.Produce(ctx, model.NewAcceptTripCommand(tripId, userId).ToKafkaMessage(), "ACCEPT")
}

func (d driverService) StartTrip(ctx context.Context, userId string, tripId string) error {
	trip, err := d.repo.GetTrip(ctx, userId, tripId)
	if err != nil {
		return err
	}
	if trip == nil {
		return helpers.ErrNotFound
	}
	err = d.repo.ChangeTripStatus(ctx, userId, tripId, StatusStarted, nil)
	if err != nil {
		return err
	}
	return d.producer.Produce(ctx, model.NewStartTripCommand(tripId).ToKafkaMessage(), "START")
}

func (d driverService) EndTrip(ctx context.Context, userId string, tripId string) error {
	trip, err := d.repo.GetTrip(ctx, userId, tripId)
	if err != nil {
		return err
	}
	if trip == nil {
		return helpers.ErrNotFound
	}
	err = d.repo.ChangeTripStatus(ctx, userId, tripId, StatusEnded, nil)
	if err != nil {
		return err
	}
	return d.producer.Produce(ctx, model.NewEndTripCommand(tripId).ToKafkaMessage(), "END")
}

func New(repo repository.Driver, producer kafka_producer.Producer, locationBaseUrl string) service.Driver {
	return &driverService{
		repo:            repo,
		producer:        producer,
		locationBaseUrl: locationBaseUrl,
	}
}
