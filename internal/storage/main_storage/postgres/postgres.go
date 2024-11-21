package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/sso/internal/config"
	"github.com/notblinkyet/sso/internal/models"
	storage "github.com/notblinkyet/sso/internal/storage/main_storage"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, host, database, username, password string, port int) (*Postgres, error) {

	const op = "storage.main_storage.postgres.NewPostgres"

	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)

	cfg, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, cfg)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{pool: pool}, nil
}

func NewPostgresFromConfig(config *config.Config) (*Postgres, error) {
	return NewPostgres(context.Background(), config.Storage.Host, config.Storage.Database,
		config.Storage.Username, os.Getenv("POSTGRES_PASS"), config.Storage.Port)
}

func (p *Postgres) SaveUser(ctx context.Context, login string, passHash []byte) (int64, error) {

	const op = "storage.main_storage.postgres.SaveUser"

	query := `
        INSERT INTO users ("login", pass_hash)
        VALUES ($1, $2)
        RETURNING id;
    `
	var id int64
	err := p.pool.QueryRow(ctx, query, login, passHash).Scan(&id)
	if err != nil {

		pgErr, ok := err.(*pgconn.PgError)

		if ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrLoginExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *Postgres) User(ctx context.Context, login string) (*models.User, error) {

	const op = "storage.main_storage.postgres.User"

	var (
		id       int64
		passHash []byte
	)

	query := `
        SELECT id, pass_hash
		FROM users
		WHERE "login"=$1
    `

	err := p.pool.QueryRow(ctx, query, login).Scan(&id, &passHash)

	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return models.NewUser(id, login, passHash), nil
}

func (p *Postgres) App(ctx context.Context, id int) (*models.App, error) {

	const op = "storage.main_storage.postgres.app"

	var (
		name   string
		secret string
	)

	query := `
		SELECT "name", "secret" FROM apps
		WHERE id = $1
	`

	err := p.pool.QueryRow(ctx, query, id).Scan(&name, &secret)

	if err != nil {

		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return models.NewApp(id, name, secret), nil
}

func (p *Postgres) IsAdmin(ctx context.Context, userID int64) (bool, error) {

	const op = "storage.main_storage.postgres.isadmin"

	var isAdmin bool

	query := `
		SELECT is_admin FROM users
		WHERE id=$1
	`

	err := p.pool.QueryRow(ctx, query, userID).Scan(&isAdmin)

	if err != nil {

		if err == pgx.ErrNoRows {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, err
}
