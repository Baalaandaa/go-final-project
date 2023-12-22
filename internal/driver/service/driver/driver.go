package driver

import (
	"final-project/internal/driver/service"
)

type driverService struct {
}

func New() service.Driver {
	return &driverService{}
}
