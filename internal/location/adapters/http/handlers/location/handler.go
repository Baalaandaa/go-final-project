package locationHandler

import (
	"encoding/json"
	"final-project/internal/location/model"
	"final-project/internal/location/service"
	"final-project/pkg/helpers"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/juju/zaputil/zapctx"
	"go.opentelemetry.io/otel"
)

type locationHandler struct {
	locationService service.Location
}

var tracer = otel.Tracer("location_service")

func (l locationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "locationHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context()).Sugar()
	logger.Info("UpdateLocationCalled")

	var location model.LatLngLiteral

	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		helpers.WriteError(w, err)
		return
	}

	var driverId = chi.URLParam(r, "driverID")

	logger.Infof("Update driver %+v location", driverId)

	err := l.locationService.UpdateLocation(r.Context(), &model.Driver{
		Lat:      location.Lat,
		Lng:      location.Lng,
		DriverId: driverId,
	})
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l locationHandler) GetNearbyDrivers(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "locationHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context()).Sugar()
	logger.Info("GetNearbyDrivers")

	var location model.LatLngLiteral

	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		helpers.WriteError(w, err)
		return
	}

	drivers, err := l.locationService.GetNearbyDrivers(r.Context(), &location)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, drivers)
}

func New(locationService service.Location) LocationHandler {
	return &locationHandler{locationService: locationService}
}
