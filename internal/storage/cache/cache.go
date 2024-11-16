package cache

import (
	"context"

	"github.com/notblinkyet/sso/internal/models"
)

type Cache interface {
	SetUser(ctx context.Context, email string, passHash []byte) error
	GetUser(ctx context.Context, email string) (models.User, error)
}
