package httpadapter

import (
	"context"
	locationHandler "final-project/internal/location/adapters/http/handlers/location"
	"final-project/internal/location/service"
	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"go.uber.org/zap"
	"moul.io/chizap"
	"net/http"
)

type adapter struct {
	config *Config

	locationHandler locationHandler.LocationHandler

	logger *zap.Logger
	server *http.Server
}

func (a adapter) Serve() error {
	r := chi.NewRouter()

	r.Use(chiprometheus.NewMiddleware("location_service"))
	r.Use(otelchi.Middleware("location_service"))
	r.Use(chizap.New(a.logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	r.Handle("/metrics", promhttp.Handler())

	apiRouter := chi.NewRouter()

	apiRouter.Get("/drivers", a.locationHandler.GetNearbyDrivers)
	apiRouter.Post("/drivers/{driverID}/location", a.locationHandler.UpdateLocation)

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

func New(ctx context.Context, config *Config, locationService service.Location) Adapter {

	return &adapter{
		config:          config,
		logger:          zapctx.Logger(ctx),
		locationHandler: locationHandler.New(locationService),
	}

}
