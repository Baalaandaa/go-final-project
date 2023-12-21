package app

import (
	httpadapter "final-project/internal/location/adapters/http"
	"time"
)

const (
	AppName                = "location"
	DefaultServeAddress    = "localhost:5355"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/location/v1"
)

type AppConfig struct {
	Debug           bool          `env:"DEBUG"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

type DatabaseConfig struct {
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig

	HTTP httpadapter.Config
}

func NewConfig() (*Config, error) {
	cnf := Config{
		App: AppConfig{
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		Database: DatabaseConfig{},
		HTTP: httpadapter.Config{
			ServeAddress: DefaultServeAddress,
			BasePath:     DefaultBasePath,
		},
	}

	return &cnf, nil
}
