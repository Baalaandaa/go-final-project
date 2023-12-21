package httpadapter

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type adapter struct {
	config *Config

	server *http.Server
}

func (a adapter) Serve() error {
	r := chi.NewRouter()

	apiRouter := chi.NewRouter()

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

func New(config *Config) Adapter {

	// TODO: swagger address

	return &adapter{
		config: config,
	}

}
