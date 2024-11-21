package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/notblinkyet/sso/internal/config"
	"github.com/notblinkyet/sso/internal/models"
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

func NewRedisFromConfig(config *config.Config) (*Redis, error) {
	return NewRedis(fmt.Sprintf("%s:%d", config.Cache.Host, config.Cache.Port), os.Getenv("REDIS_PASS"), config.Cache.DB)
}

func (r *Redis) SetUser(ctx context.Context, login string, passHash []byte, id int64, expiration time.Duration) error {

	const op = "storage.cache.redis.SetUser"

	user := models.NewUser(id, login, passHash)

	jsonData, err := json.Marshal(user)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = r.client.Set(ctx, login, jsonData, expiration).Err()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) GetUser(ctx context.Context, login string) (*models.User, error) {

	var user models.User

	const op = "storage.cache.redis.GetUser"

	res, err := r.client.Get(ctx, login).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("%s: %w", op, cache.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = json.Unmarshal([]byte(res), &user)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
