package repository

import "context"

type Location interface {
	UpdateLocation(ctx context.Context, lat, lng float64, driverID string) error
	GetNearbyDrivers(ctx context.Context, lat, lng, radius float64) error
}
