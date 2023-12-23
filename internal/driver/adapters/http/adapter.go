package httpadapter

import (
	"context"
	driverHandler "final-project/internal/driver/adapters/http/handlers/driver"
	"final-project/internal/driver/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juju/zaputil/zapctx"
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

	//r.Use(otelchi.Middleware("driver_service"))
	r.Use(chizap.New(a.logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	apiRouter := chi.NewRouter()

	// apiRouter.Get("/some_method", http.HandlerFunc(a.driverHandler.SomeMethod))

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

	// TODO: swagger address

	return &adapter{
		config:        config,
		logger:        zapctx.Logger(ctx),
		driverHandler: driverHandler.New(driverService),
	}

}