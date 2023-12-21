package httpadapter

type Config struct {
	ServeAddress string `env:"SERVE_ADDRESS"`
	BasePath     string `env:"BASE_PATH"`

	UseTLS     bool   `env:"USE_TLS"`
	TLSKeyFile string `env:"TLS_KEY_FILE"`
	TLSCrtFile string `env:"TLS_CRT_FILE"`

	SwaggerAddress string `env:"SWAGGER_ADDRESS"`
}
