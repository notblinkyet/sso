package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string, password string, db int) *Redis {

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Redis{
		client: client,
	}
}

func (r *Redis) SetUser(ctx context.Context, email string, passHash []byte) error {
	status := r.client.Set(ctx, email, string(passHash), time.Duration(time.Hour))

	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (r *Redis) GetUser(ctx context.Context, email string) ([]byte, error) {
	res, err := r.client.Get(ctx, email).Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}
