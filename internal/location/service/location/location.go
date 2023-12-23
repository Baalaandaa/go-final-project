package location

import (
	"context"
	"final-project/internal/location/model"
	"final-project/internal/location/repository"
	"final-project/internal/location/service"
)

type locationService struct {
	repo repository.Location
}

func (l locationService) UpdateLocation(ctx context.Context, driver *model.Driver) error {
	return l.repo.UpdateLocation(ctx, driver.Lat, driver.Lng, driver.DriverId)
}

func (l locationService) GetNearbyDrivers(ctx context.Context, location *model.LatLngLiteral) ([]*model.Driver, error) {
	return l.repo.GetNearbyDrivers(ctx, location.Lat, location.Lng, location.Radius)
}

func New(repo repository.Location) service.Location {
	return &locationService{
		repo: repo,
	}
}
