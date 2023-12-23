package repository

import (
	"context"
	"final-project/internal/location/model"
)

type Location interface {
	UpdateLocation(ctx context.Context, lat, lng float64, driverID string) error
	GetNearbyDrivers(ctx context.Context, lat, lng, radius float64) ([]*model.Driver, error)
}
