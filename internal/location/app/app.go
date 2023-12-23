package app

import (
	"context"
	"final-project/internal/location/adapters/http"
	location_repo "final-project/internal/location/repository/location"
	"final-project/internal/location/service"
	"final-project/internal/location/service/location"
	"final-project/pkg/helpers"
	"final-project/pkg/logger"
	"final-project/pkg/otel"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	config *Config
	logger *zap.Logger

	locationService service.Location

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

	pgxPool, err := helpers.InitPostgres(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.Postgres.User, config.Postgres.Password, config.Postgres.Host, config.Postgres.Database), config.Postgres.MigrationsPath)
	if err != nil {
		l.Sugar().Fatalf("InitPostgres err: %+v", err)
		return nil, err
	}

	locationRepo := location_repo.New(pgxPool)
	locationService := location.New(locationRepo)

	return &app{
		config:          config,
		logger:          l,
		locationService: locationService,
		httpAdapter:     httpadapter.New(ctx, &config.HTTP, locationService),
	}, nil
}
