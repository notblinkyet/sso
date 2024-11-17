package main

import (
	"log/slog"
	"os"

	"github.com/notblinkyet/sso/internal/config"
	"github.com/notblinkyet/sso/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {

	// TODO: Read configuration from yaml file

	config := config.MustLoad()

	// TODO: Init logger

	logger := setupLogger(config.Env)

	logger.Info("", slog.Any("config", config))

	logger.Info("Application started")

	// TODO: Create application

	// TODO: Start application

	// TODO: Graceful shutdown

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
