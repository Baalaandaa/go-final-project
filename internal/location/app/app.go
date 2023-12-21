package app

import (
	"context"
	"final-project/internal/location/adapters/http"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type app struct {
	config *Config

	httpAdapter httpadapter.Adapter
}

func (a app) Serve() error {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.httpAdapter.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	<-done

	a.Shutdown()

	return nil
}

func (a app) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
	defer cancel()

	a.httpAdapter.Shutdown(ctx)
}

func New(config *Config) (App, error) {
	return &app{
		config:      config,
		httpAdapter: httpadapter.New(&config.HTTP),
	}, nil
}
