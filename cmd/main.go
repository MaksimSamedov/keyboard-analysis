package main

import (
	"keyboard-analysis/internal/app"
	"keyboard-analysis/internal/config"
	"log"
)

func main() {
	conf := config.WithDefaults()
	logger := log.Default()
	application := app.New(conf, logger)

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
