package httpadapter

import (
	"context"
	locationHandler "final-project/internal/location/adapters/http/handlers/location"
	"final-project/internal/location/service"
	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
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

	//r.Use(otelchi.Middleware("location_service"))
	r.Use(chizap.New(a.logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	apiRouter := chi.NewRouter()

	apiRouter.Get("/drivers", http.HandlerFunc(a.locationHandler.GetNearbyDrivers))
	apiRouter.Post("/drivers/{driverID}/location", http.HandlerFunc(a.locationHandler.UpdateLocation))

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

	// TODO: swagger address

	return &adapter{
		config:          config,
		logger:          zapctx.Logger(ctx),
		locationHandler: locationHandler.New(locationService),
	}

}
