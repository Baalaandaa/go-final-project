package service

import (
	"context"
	"final-project/internal/location/model"
)

type Location interface {
	UpdateLocation(ctx context.Context, driver *model.Driver) error
	GetNearbyDrivers(ctx context.Context, location *model.LatLngLiteral) ([]*model.Driver, error)
}
