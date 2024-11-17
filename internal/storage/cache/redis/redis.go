package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/notblinkyet/sso/internal/storage/cache"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string, password string, db int) (*Redis, error) {

	const op = "storage.cache.redis.NewRedis"

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) SetUser(ctx context.Context, login string, passHash []byte, expiration time.Duration) error {

	const op = "storage.cache.redis.SetUser"

	status := r.client.Set(ctx, login, string(passHash), expiration)

	if status.Err() != nil {
		return fmt.Errorf("%s: %w", op, status.Err())
	}
	return nil
}

func (r *Redis) GetUser(ctx context.Context, login string) ([]byte, error) {

	const op = "storage.cache.redis.GetUser"

	res, err := r.client.Get(ctx, login).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("%s: %w", op, cache.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return []byte(res), nil
}
