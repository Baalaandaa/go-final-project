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
	return nil
}

func (l locationService) GetNearbyDrivers(ctx context.Context, location *model.LatLngLiteral) ([]*model.Driver, error) {
	drivers := []*model.Driver{
		&model.Driver{
			Lat:      0.5,
			Lng:      0.4,
			DriverId: "abacaba",
		},
	}
	return drivers, nil
}

func New(repo repository.Location) service.Location {
	return &locationService{
		repo: repo,
	}
}
