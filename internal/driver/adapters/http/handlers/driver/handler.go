package driverHandler

import (
	"final-project/internal/driver/service"
	"final-project/pkg/helpers"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
	"go.opentelemetry.io/otel"
)

type driverHandler struct {
	driverService service.Driver
}

var tracer = otel.Tracer("location_service")

func (dh *driverHandler) ListTrips(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "driverHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context())
	logger.Info("ListTripsCalled")

	userId := r.Header.Get("user_id")
	trips, err := dh.driverService.ListTrips(r.Context(), userId)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, trips)
}

func (dh *driverHandler) GetTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "driverHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context())
	logger.Info("GetTripCalled")

	userId := r.Header.Get("user_id")
	tripId := chi.URLParam(r, "trip_id")

	trip, err := dh.driverService.GetTrip(r.Context(), userId, tripId)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, trip)
}

func (dh *driverHandler) CancelTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "driverHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context())
	logger.Info("CancelTripCalled")

	userId := r.Header.Get("user_id")
	tripId := chi.URLParam(r, "trip_id")
	reason := r.URL.Query().Get("reason")

	err := dh.driverService.CancelTrip(r.Context(), userId, tripId, &reason)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (dh *driverHandler) AcceptTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "driverHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context())
	logger.Info("AcceptTripCalled")

	userId := r.Header.Get("user_id")
	tripId := chi.URLParam(r, "trip_id")

	err := dh.driverService.AcceptTrip(r.Context(), userId, tripId)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (dh *driverHandler) StartTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "driverHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context())
	logger.Info("StartTripCalled")

	userId := r.Header.Get("user_id")
	tripId := chi.URLParam(r, "trip_id")

	err := dh.driverService.StartTrip(r.Context(), userId, tripId)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (dh *driverHandler) EndTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "driverHandler")
	defer span.End()

	logger := zapctx.Logger(r.Context())
	logger.Info("EndTripCalled")

	userId := r.Header.Get("user_id")
	tripId := chi.URLParam(r, "trip_id")

	err := dh.driverService.EndTrip(r.Context(), userId, tripId)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func New(driverService service.Driver) DriverHandler {
	return &driverHandler{driverService: driverService}
}
