package cache

import (
	"context"
	"errors"
	"time"

	"github.com/notblinkyet/sso/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Cache interface {
	SetUser(ctx context.Context, login string, passHash []byte, id int64, expiration time.Duration) error
	GetUser(ctx context.Context, login string) (*models.User, error)
}
