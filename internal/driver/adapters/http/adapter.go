package httpadapter

import (
	"context"
	driverHandler "final-project/internal/driver/adapters/http/handlers/driver"
	"final-project/internal/driver/service"
	chiprometheus "github.com/766b/chi-prometheus"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
	"github.com/riandyrn/otelchi"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"
	"moul.io/chizap"
)

type adapter struct {
	config *Config

	driverHandler driverHandler.DriverHandler

	logger *zap.Logger
	server *http.Server
}

func (a adapter) Serve() error {
	r := chi.NewRouter()

	r.Use(chiprometheus.NewMiddleware("location_service"))
	r.Use(otelchi.Middleware("driver_service"))
	r.Use(chizap.New(a.logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	r.Handle("/metrics", promhttp.Handler())

	apiRouter := chi.NewRouter()

	apiRouter.Get("/trips", a.driverHandler.ListTrips)
	apiRouter.Get("/trips/{trip_id}", a.driverHandler.GetTrip)
	apiRouter.Post("/trips/{trip_id}/cancel", a.driverHandler.CancelTrip)
	apiRouter.Post("/trips/{trip_id}/accept", a.driverHandler.AcceptTrip)
	apiRouter.Post("/trips/{trip_id}/start", a.driverHandler.StartTrip)
	apiRouter.Post("/trips/{trip_id}/end", a.driverHandler.EndTrip)

	r.Mount(a.config.BasePath, apiRouter)
	a.server = &http.Server{Addr: a.config.ServeAddress, Handler: r}

	if a.config.UseTLS {
		return a.server.ListenAndServeTLS(a.config.TLSCrtFile, a.config.TLSKeyFile)
	}

	return a.server.ListenAndServe()
}

func (a adapter) Shutdown(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
}

func New(ctx context.Context, config *Config, driverService service.Driver) Adapter {

	return &adapter{
		config:        config,
		logger:        zapctx.Logger(ctx),
		driverHandler: driverHandler.New(driverService),
	}

}
