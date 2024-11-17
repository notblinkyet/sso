package storage

import (
	"context"
	"errors"

	models "github.com/notblinkyet/sso/internal/models"
)

var (
	ErrLoginExists = errors.New("login already exists")
)

type Storage interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (int64, error)
	User(ctx context.Context, login string) (models.User, error)
	App(ctx context.Context, id int) (models.App, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
