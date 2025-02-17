package main

import (
	"os"

	"github.com/givko/hoodie/internal/api/router"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"
	zerologr.SetMaxV(1)

	var zl zerolog.Logger
	if os.Getenv("ENV") == "dev" {
		zl = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05", // Customize time format if needed.
		}).With().Timestamp().Logger()
	} else {
		// In production, use the default JSON logger.
		zl = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	var log logr.Logger = zerologr.New(&zl)
	router.Init(log).Run(":8080")
}
