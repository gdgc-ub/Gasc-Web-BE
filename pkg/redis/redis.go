package redis

import (
	"context"
	"errors"
	redisPkg "github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

type ItfRedis interface {
	SetOTP(ctx context.Context, userID string, code string) error
	GetOTP(ctx context.Context, userID string) (string, error)
}

type redis struct {
	client *redisPkg.Client
}

func New() ItfRedis {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redisPkg.NewClient(&redisPkg.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return &redis{client: client}
}

func (r *redis) SetOTP(ctx context.Context, userEmail string, code string) error {
	err := r.client.Set(ctx, code, userEmail, 2*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redis) GetOTP(ctx context.Context, code string) (string, error) {
	val, err := r.client.Get(ctx, code).Result()
	if errors.Is(err, redisPkg.Nil) {
		// Key does not exist
		return "", nil
	} else if err != nil {

		return "", err
	}
	return val, nil
}
