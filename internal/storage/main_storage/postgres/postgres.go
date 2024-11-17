package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/sso/internal/models"
	storage "github.com/notblinkyet/sso/internal/storage/main_storage"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, host, port, database, username, password string) (*Postgres, error) {

	const op = "storage.main_storage.postgres.NewPostgres"

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)

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

func (p *Postgres) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {

	const op = "storage.main_storage.postgres.SaveUser"

	query := `
        INSERT INTO users (email, pass_hash)
        VALUES ($1, $2)
        RETURNING id;
    `
	var id int64
	err := p.pool.QueryRow(ctx, query, email, passHash).Scan(&id)
	if err != nil {
		var pgErr pgx.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *Postgres) User(ctx context.Context, email string) (models.User, error) {

	const op = "storage.main_storage.postgres.User"

	var (
		id       int64
		passHash []byte
	)

	query := `
        SELECT id, pass_hash
		FROM users
		WHERE email=$1
    `

	err := p.pool.QueryRow(ctx, query, email).Scan(&id, &passHash)

	if err != nil {

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return *models.NewUser(id, email, passHash), nil
}

func (p *Postgres) App(ctx context.Context, id int) (models.App, error) {

	const op = "storage.main_storage.postgres.app"

	var (
		name   string
		secret string
	)

	query := `
		SELECT name, secret FROM app
		WHERE id = $1
	`

	err := p.pool.QueryRow(ctx, query, id).Scan(&name, &secret)

	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return *models.NewApp(id, name, secret), nil
}

func (p *Postgres) IsAdmin(ctx context.Context, userID int64) (bool, error) {

	const op = "storage.main_storage.postgres.isadmin"

	var isAdmin bool

	query := `
		SELECT is_admin FROM user
		WHERE id=$1
	`

	err := p.pool.QueryRow(ctx, query, userID).Scan(&isAdmin)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, err
}
