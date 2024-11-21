package app

import (
	"log/slog"

	grpcapp "github.com/notblinkyet/sso/internal/app/grpc"
	"github.com/notblinkyet/sso/internal/config"
	"github.com/notblinkyet/sso/internal/lib/logger/sl"
	"github.com/notblinkyet/sso/internal/services/auth"
	"github.com/notblinkyet/sso/internal/storage/cache"
	"github.com/notblinkyet/sso/internal/storage/cache/redis"
	storage "github.com/notblinkyet/sso/internal/storage/main_storage"
	"github.com/notblinkyet/sso/internal/storage/main_storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, config *config.Config) *App {

	var (
		main_storage storage.Storage
		cache        cache.Cache
		err          error
	)

	switch config.Storage.Type {
	case "postgres":
		main_storage, err = postgres.NewPostgresFromConfig(config)

		if err != nil {
			log.Error("failed to connect to Posgres", sl.Err(err))
			panic(err)
		}
	default:
		panic("Unknown DB")
	}

	switch config.Cache.Driver {
	case "redis":
		cache, err = redis.NewRedisFromConfig(config)
		if err != nil {
			log.Error("failed to connect to Redis", sl.Err(err))
			panic(err)
		}
	default:
		panic("Unknown Cache")
	}

	auth := auth.New(log, main_storage, cache, config.TokenTTL)

	GRPCServer := grpcapp.New(log, auth, config.Grpc.Port, config.Grpc.Timeout)

	return &App{
		GRPCServer: GRPCServer,
	}
}
