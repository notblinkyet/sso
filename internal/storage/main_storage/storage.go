package storage

import (
	"context"
	"errors"

	models "github.com/notblinkyet/sso/internal/models"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type Storage interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
	User(ctx context.Context, email string) (models.User, error)
	App(ctx context.Context, id int) (models.App, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
