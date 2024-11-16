package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/notblinkyet/sso/internal/models"
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgres(ctx context.Context, host, port, database, username, password string) (*Postgres, error) {

	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, database, username, password)

	conn, err := pgx.Connect(ctx, connStr)

	if err != nil {
		return nil, err
	}
	return &Postgres{conn: conn}, nil
}

func (p *Postgres) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	query := `
        INSERT INTO users (email, pass_hash)
        VALUES ($1, $2)
        RETURNING id;
    `
	var id int64
	err := p.conn.QueryRow(ctx, query, email, passHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *Postgres) User(ctx context.Context, email string) (models.User, error) {
	var (
		id       int64
		passHash []byte
	)

	query := `
        SELECT id, pass_hash
		FROM users
		WHERE email=$1
    `

	err := p.conn.QueryRow(ctx, query, email).Scan(id, passHash)

	return models.User{
		Id:       id,
		Email:    email,
		PassHash: passHash,
	}, err
}

func (p *Postgres) App(ctx context.Context, id int) (models.App, error) {

	var (
		name   string
		secret string
	)

	query := `
		SELECT name, secret FROM app
		WHERE id = $1
	`

	err := p.conn.QueryRow(ctx, query, id).Scan(name, secret)

	return models.App{
		ID:     id,
		Name:   name,
		Secret: secret,
	}, err
}

func (p *Postgres) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	var IsAdmin bool

	query := `
		SELECT isadmin FROM user
		WHERE id=$1
	`

	err := p.conn.QueryRow(ctx, query, userID).Scan(&IsAdmin)

	return IsAdmin, err
}