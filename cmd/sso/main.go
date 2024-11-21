package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/notblinkyet/sso/internal/app"
	"github.com/notblinkyet/sso/internal/config"
	"github.com/notblinkyet/sso/internal/logger"
)

func main() {

	// TODO: Read configuration from yaml file

	config := config.MustLoad()

	// TODO: Init logger

	logger := logger.SetupLogger(config.Env)

	logger.Info("", slog.Any("config", config))

	logger.Info("Success read config and setup logger")

	// TODO: Create application

	app := app.New(logger, config)

	go func() {
		app.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	logger.Info("Gracefull stopped")

}
