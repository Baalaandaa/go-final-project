package app

import (
	httpadapter "final-project/internal/location/adapters/http"
	config_reader "final-project/pkg/config-reader"
)

const (
	AppName                = "location_service"
	DefaultServeAddress    = "localhost:5355"
	DefaultShutdownTimeout = 20
	DefaultBasePath        = "/location/v1"
	DefaultOTLPEndpoint    = "0.0.0.0:4317"
)

type AppConfig struct {
	Debug                    bool `env:"DEBUG"`
	ShutdownTimeoutInSeconds int  `env:"SHUTDOWN_TIMEOUT"`
}

type DatabaseConfig struct {
}

type OTLPConfig struct {
	Endpoint string `env:"OTLP_ENDPOINT"`
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	OTLP     OTLPConfig

	HTTP httpadapter.Config
}

func NewConfig() (*Config, error) {
	cnf := Config{
		App: AppConfig{
			ShutdownTimeoutInSeconds: DefaultShutdownTimeout,
		},
		Database: DatabaseConfig{},
		HTTP: httpadapter.Config{
			ServeAddress: DefaultServeAddress,
			BasePath:     DefaultBasePath,
		},
		OTLP: OTLPConfig{
			Endpoint: DefaultOTLPEndpoint,
		},
	}

	config_reader.ReadEnv(&cnf)

	return &cnf, nil
}
