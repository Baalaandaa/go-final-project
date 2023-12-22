package main

import (
	"context"
	"final-project/internal/driver/app"
	"log"
)

func main() {
	config, err := app.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	application, err := app.New(context.Background(), config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := application.Serve(); err != nil {
		log.Fatal(err.Error())
	}
}
