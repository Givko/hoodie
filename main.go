package main

import (
	"os"

	"github.com/givko/hoodie/internal/api/router"
	"github.com/givko/hoodie/internal/api/ws"
	"github.com/givko/hoodie/internal/infrastructure/subscribers"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var Hub ws.WsHubInterface

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"
	zerologr.SetMaxV(1)

	var zl zerolog.Logger
	if os.Getenv("ENV") == "dev" {
		zl = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Logger()
	} else {
		//Log into file in order for the logs to be persisted\
		file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		zl = zerolog.New(file).With().Timestamp().Logger()
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	var log logr.Logger = zerologr.New(&zl)

	Hub = ws.NewHub(log, redisClient)
	go subscribers.StartChatMessageSubscriber(redisClient, log)
	go Hub.Run()

	router.Init(Hub, log, redisClient).Run(":8080")

}
