package services

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenBlacklist struct {
	client *redis.Client
}

func NewTokenBlacklist(client *redis.Client) *TokenBlacklist {
	return &TokenBlacklist{client: client}
}

func (t *TokenBlacklist) Add(token string, expiry time.Time) error {
	ctx := context.Background()
	ttl := time.Until(expiry)
	return t.client.Set(ctx, "invalidated:"+token, "1", ttl).Err()
}

func (t *TokenBlacklist) Exists(token string) (bool, error) {
	ctx := context.Background()
	val, err := t.client.Exists(ctx, "invalidated:"+token).Result()
	return val > 0, err
}
