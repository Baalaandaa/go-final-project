package app

import (
	"context"
	httpadapter "final-project/internal/driver/adapters/http"
	driver_repo "final-project/internal/driver/repository/driver"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ConnectMongoDB(uri string, name string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(name)
	return db, nil
}

func New(ctx context.Context, config *Config) (App, error) {
	l, err := logger.GetLogger(config.App.Debug)
	if err != nil {
		return nil, err
	}

	ctx = zapctx.WithLogger(ctx, l)

	db, err := ConnectMongoDB(config.Database.DatabaseUri, config.Database.DatabaseName)
	if err != nil {
		return nil, err
	}

	driverRepo := driver_repo.New(db, config.Database.DatabaseName)
	driverService := driver.New(driverRepo)

	return &app{
		config:        config,
		logger:        l,
		driverService: driverService,
		httpAdapter:   httpadapter.New(ctx, &config.HTTP, driverService),
	}, nil
}
