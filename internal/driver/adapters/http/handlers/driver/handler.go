package driverHandler

import (
	"final-project/internal/driver/service"

	"go.opentelemetry.io/otel"
)

type driverHandler struct {
	driverService service.Driver
}

var tracer = otel.Tracer("location_service")

func New(driverService service.Driver) DriverHandler {
	return &driverHandler{driverService: driverService}
}
