package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/notblinkyet/sso/internal/config"
	"github.com/notblinkyet/sso/internal/logger"
)

func main() {
	var rollbackSteps int

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Storage.Username,
		os.Getenv("POSTGRES_PASS"), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database)
	migrationsPath := cfg.MigrationsPath

	flag.IntVar(&rollbackSteps, "rollback-steps", 0, "number of steps to rollback (use negative value for specific behavior)")
	flag.Parse()

	if dbURL == "" {
		log.Error("db-url is required")
		panic("db-url is required")
	}
	if migrationsPath == "" {
		log.Error("migrations-path is required")
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("%s&x-migrations-table=", dbURL),
	)
	if err != nil {
		log.Error("failed to create migration engine: %v", slog.String("err", err.Error()))
		panic(err)
	}

	if rollbackSteps > 0 {
		if err := m.Steps(-rollbackSteps); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to rollback")
				return
			}
			log.Error("failed to migrate rollback: %v", slog.String("err", err.Error()))
			panic(err)
		}
		log.Info("rolled back steps\n", slog.Int("", rollbackSteps))
		return
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Error("failed to migrate up: %v", slog.String("err", err.Error()))
		panic(err)
	}

	log.Info("migrations applied successfully")
}
