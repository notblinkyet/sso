package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/notblinkyet/sso/internal/lib/jwt"
	"github.com/notblinkyet/sso/internal/lib/logger/sl"
	"github.com/notblinkyet/sso/internal/storage/cache"
	storage "github.com/notblinkyet/sso/internal/storage/main_storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log      *slog.Logger
	storage  storage.Storage
	cache    cache.Cache
	tokenTTL time.Duration
}

func New(log *slog.Logger, storage storage.Storage, cache cache.Cache, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storage,
		cache:    cache,
		tokenTTL: tokenTTL,
	}
}

func (auth *Auth) Register(ctx context.Context, login, password string) (int64, error) {

	const op = "services.auth.register"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate passhash", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := auth.storage.SaveUser(ctx, login, passHash)

	if err != nil {
		if errors.Is(err, storage.ErrLoginExists) {

			log.Warn("this login already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, storage.ErrLoginExists)
		}

		log.Error("failed to save user in main storage", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = auth.cache.SetUser(ctx, login, passHash, id, 24*time.Hour)

	if err != nil {
		log.Error("failed to save user in cache", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user register successfully")

	return id, nil
}

func (auth *Auth) Login(ctx context.Context, login string, password string, appID int) (string, error) {

	const op = "services.auth.login"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("login user")

	user, err := auth.cache.GetUser(ctx, login)

	if err != nil {

		if !errors.Is(err, cache.ErrUserNotFound) {

			log.Error("failed with get data from cache", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)

		}

		user, err = auth.storage.User(ctx, login)

		if err != nil {

			if errors.Is(err, storage.ErrUserNotFound) {

				log.Warn("user not found", sl.Err(storage.ErrUserNotFound))

				return "", fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)

			}

			log.Error("failed with get data from main storage", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))

	if err != nil {
		log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := auth.storage.App(ctx, appID)

	if err != nil {
		log.Error("failed to get app", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	err = auth.cache.SetUser(ctx, user.Login, user.PassHash, user.ID, 24*time.Hour)

	if err != nil {

		log.Error("failed to save user in cache", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, auth.tokenTTL)

	if err != nil {

		log.Error("failed to generate tocken", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)

	}

	return token, err
}

func (auth *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {

	const op = "services.auth.is_admin"

	log := auth.log.With(
		slog.String("op", op),
		slog.Int64("id", userID),
	)

	log.Info("get info about rules")

	isAdmin, err := auth.storage.IsAdmin(ctx, userID)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (auth *Auth) Logout(ctx context.Context, token string) (bool, error) {
	// const op = "services.auth.logout"
	panic("not implement")
}
