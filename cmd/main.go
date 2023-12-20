package main

import (
	"keyboard-analysis/internal/app"
	"keyboard-analysis/internal/config"
	"log"
)

func main() {
	// configure
	conf := config.WithDefaults()
	conf.PasswordsCount = 3 // 4
	conf.AnalyserProps.MaxDeviation = 50
	conf.AnalyserProps.MinSuccessfulComparisons = 50

	// prepare app
	logger := log.Default()
	application := app.New(conf, logger)

	// run
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
