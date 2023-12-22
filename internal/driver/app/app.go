package app

import (
	"context"
	"final-project/internal/driver/adapters/http"
	"final-project/internal/driver/service"
	"final-project/internal/driver/service/driver"
	"final-project/pkg/logger"
	"final-project/pkg/otel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type app struct {
	config *Config
	logger *zap.Logger

	driverService service.Driver

	httpAdapter httpadapter.Adapter
}

func (a app) Serve() error {
	shutdown := otel.InitProvider(AppName, a.config.OTLP.Endpoint)
	defer shutdown()

	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("Start serving")
		if err := a.httpAdapter.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Serve error", err.Error())
		}
	}()

	<-done

	a.Shutdown()

	return nil
}

func (a app) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.App.ShutdownTimeoutInSeconds)*time.Second)
	defer cancel()

	a.httpAdapter.Shutdown(ctx)
}

func New(ctx context.Context, config *Config) (App, error) {
	l, err := logger.GetLogger(config.App.Debug)
	if err != nil {
		return nil, err
	}

	ctx = zapctx.WithLogger(ctx, l)

	driverService := driver.New()

	return &app{
		config:        config,
		logger:        l,
		driverService: driverService,
		httpAdapter:   httpadapter.New(ctx, &config.HTTP, driverService),
	}, nil
}
