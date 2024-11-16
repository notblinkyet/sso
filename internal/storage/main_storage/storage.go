package storage

import (
	"context"

	models "github.com/notblinkyet/sso/internal/models"
	"github.com/notblinkyet/sso/internal/storage/main_storage/postgres"
)

type Storage interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
	User(ctx context.Context, email string) (models.User, error)
	App(ctx context.Context, id int) (models.App, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func main() {
	var s Storage
	s, err := postgres.NewPostgres(context.Background(), "", "", "", "", "")
	_ = s
	_ = err
}
